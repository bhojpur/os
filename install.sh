#!/bin/bash
set -e

PROG=$0
PROGS="dd curl mkfs.ext4 mkfs.vfat fatlabel parted partprobe grub-install"
DISTRO=/run/bos/iso

if [ "$BOS_DEBUG" = true ]; then
    set -x
fi

get_url()
{
    FROM=$1
    TO=$2
    case $FROM in
        ftp*|http*|tftp*)
            n=0
            attempts=5
            until [ "$n" -ge "$attempts" ]
            do
                curl -o $TO -fL ${FROM} && break
                n=$((n+1))
                echo "Failed to download, retry attempt ${n} out of ${attempts}"
                sleep 2
            done
            ;;
        *)
            cp -f $FROM $TO
            ;;
    esac
}

cleanup2()
{
    if [ -n "${TARGET}" ]; then
        umount ${TARGET}/boot/efi || true
        umount ${TARGET} || true
    fi

    losetup -d ${ISO_DEVICE} || losetup -d ${ISO_DEVICE%?} || true
    umount $DISTRO || true
}

cleanup()
{
    EXIT=$?
    cleanup2 2>/dev/null || true
    return $EXIT
}

usage()
{
    echo "Usage: $PROG [--force-efi] [--debug] [--tty TTY] [--poweroff] [--takeover] [--no-format] [--config https://.../config.yaml] DEVICE ISO_URL"
    echo ""
    echo "Example: $PROG /dev/vda https://github.com/bhojpur/os/releases/download/v0.0.2/bos.iso"
    echo ""
    echo "DEVICE must be the disk that will be partitioned (/dev/vda). If you are using --no-format it should be the device of the BOS_STATE partition (/dev/vda2)"
    echo ""
    echo "The parameters names refer to the same names used in the command line, refer to README.md for"
    echo "more info."
    echo ""
    exit 1
}

do_format()
{
    if [ "$BOS_INSTALL_NO_FORMAT" = "true" ]; then
        STATE=$(blkid -L BOS_STATE || true)
        if [ -z "$STATE" ] && [ -n "$DEVICE" ]; then
            tune2fs -L BOS_STATE $DEVICE
            STATE=$(blkid -L BOS_STATE)
        fi

        return 0
    fi

    dd if=/dev/zero of=${DEVICE} bs=1M count=1
    parted -s ${DEVICE} mklabel ${PARTTABLE}
    if [ "$PARTTABLE" = "gpt" ]; then
        BOOT_NUM=1
        STATE_NUM=2
        parted -s ${DEVICE} mkpart primary fat32 0% 50MB
        parted -s ${DEVICE} mkpart primary ext4 50MB 750MB
    else
        BOOT_NUM=
        STATE_NUM=1
        parted -s ${DEVICE} mkpart primary ext4 0% 700MB
    fi
    parted -s ${DEVICE} set 1 ${BOOTFLAG} on
    partprobe ${DEVICE} 2>/dev/null || true
    sleep 2

    PREFIX=${DEVICE}
    if [ ! -e ${PREFIX}${STATE_NUM} ]; then
        PREFIX=${DEVICE}p
    fi

    if [ ! -e ${PREFIX}${STATE_NUM} ]; then
        echo Failed to find ${PREFIX}${STATE_NUM} or ${DEVICE}${STATE_NUM} to format
        exit 1
    fi

    if [ -n "${BOOT_NUM}" ]; then
        BOOT=${PREFIX}${BOOT_NUM}
    fi
    STATE=${PREFIX}${STATE_NUM}

    mkfs.ext4 -F -L BOS_STATE ${STATE}
    if [ -n "${BOOT}" ]; then
        mkfs.vfat -F 32 ${BOOT}
        fatlabel ${BOOT} BOS_GRUB
    fi
}

do_mount()
{
    TARGET=/run/bos/target
    mkdir -p ${TARGET}
    mount ${STATE} ${TARGET}
    mkdir -p ${TARGET}/boot
    if [ -n "${BOOT}" ]; then
        mkdir -p ${TARGET}/boot/efi
        mount ${BOOT} ${TARGET}/boot/efi
    fi

    mkdir -p ${DISTRO}
    mount -o ro ${ISO_DEVICE} ${DISTRO} || mount -o ro ${ISO_DEVICE%?} ${DISTRO}
}

do_copy()
{
    tar cf - -C ${DISTRO} bos | tar xvf - -C ${TARGET}
    if [ -n "$STATE_NUM" ]; then
        echo $DEVICE $STATE_NUM > $TARGET/bhojpur/system/growpart
    fi

    if [ -n "$BOS_INSTALL_CONFIG_URL" ]; then
        get_url "$BOS_INSTALL_CONFIG_URL" ${TARGET}/bhojpur/system/config.yaml
        chmod 600 ${TARGET}/bhojpur/system/config.yaml
    fi

    if [ "$BOS_INSTALL_TAKE_OVER" = "true" ]; then
        touch ${TARGET}/bhojpur/system/takeover

        if [ "$BOS_INSTALL_POWER_OFF" = true ] || grep -q 'bos.install.power_off=true' /proc/cmdline; then
            touch ${TARGET}/bhojpur/system/poweroff
        fi
    fi
}

install_grub()
{
    if [ "$BOS_INSTALL_DEBUG" ]; then
        GRUB_DEBUG="bos.debug"
    fi

    mkdir -p ${TARGET}/boot/grub
    cat > ${TARGET}/boot/grub/grub.cfg << EOF
set default=0
set timeout=10

set gfxmode=auto
set gfxpayload=keep
insmod all_video
insmod gfxterm

menuentry "Bhojpur OS Current" {
  search.fs_label BOS_STATE root
  set sqfile=/bhojpur/system/kernel/current/kernel.squashfs
  loopback loop0 /\$sqfile
  set root=(\$root)
  linux (loop0)/vmlinuz printk.devkmsg=on console=tty1 $GRUB_DEBUG
  initrd /bhojpur/system/kernel/current/initrd
}

menuentry "Bhojpur OS Previous" {
  search.fs_label BOS_STATE root
  set sqfile=/bhojpur/system/kernel/previous/kernel.squashfs
  loopback loop0 /\$sqfile
  set root=(\$root)
  linux (loop0)/vmlinuz printk.devkmsg=on console=tty1 $GRUB_DEBUG
  initrd /bhojpur/system/kernel/previous/initrd
}

menuentry "Bhojpur OS Rescue (current)" {
  search.fs_label BOS_STATE root
  set sqfile=/bhojpur/system/kernel/current/kernel.squashfs
  loopback loop0 /\$sqfile
  set root=(\$root)
  linux (loop0)/vmlinuz printk.devkmsg=on rescue console=tty1
  initrd /bhojpur/system/kernel/current/initrd
}

menuentry "Bhojpur OS Rescue (previous)" {
  search.fs_label BOS_STATE root
  set sqfile=/bhojpur/system/kernel/previous/kernel.squashfs
  loopback loop0 /\$sqfile
  set root=(\$root)
  linux (loop0)/vmlinuz printk.devkmsg=on rescue console=tty1
  initrd /bhojpur/system/kernel/previous/initrd
}
EOF
    if [ -z "${BOS_INSTALL_TTY}" ]; then
        TTY=$(tty | sed 's!/dev/!!')
    else
        TTY=$BOS_INSTALL_TTY
    fi
    if [ -e "/dev/${TTY%,*}" ] && [ "$TTY" != tty1 ] && [ "$TTY" != console ] && [ -n "$TTY" ]; then
        sed -i "s!console=tty1!console=tty1 console=${TTY}!g" ${TARGET}/boot/grub/grub.cfg
    fi

    if [ "$BOS_INSTALL_NO_FORMAT" = "true" ]; then
        return 0
    fi

    if [ "$BOS_INSTALL_FORCE_EFI" = "true" ]; then
        if [ $(uname -m) = "aarch64" ]; then
            GRUB_TARGET="--target=arm64-efi"
        else
            GRUB_TARGET="--target=x86_64-efi"
        fi
    fi

    grub-install ${GRUB_TARGET} --boot-directory=${TARGET}/boot --removable ${DEVICE}
}

get_iso()
{
    ISO_DEVICE=$(blkid -L BOS || true)
    if [ -z "${ISO_DEVICE}" ]; then
        for i in $(lsblk -o NAME,TYPE -n | grep -w disk | awk '{print $1}'); do
            mkdir -p ${DISTRO}
            if mount -o ro /dev/$i ${DISTRO}; then
                ISO_DEVICE="/dev/$i"
                umount ${DISTRO}
                break
            fi
        done
    fi

    if [ -z "${ISO_DEVICE}" ] && [ -n "$BOS_INSTALL_ISO_URL" ]; then
        TEMP_FILE=$(mktemp bos.XXXXXXXX.iso)
        get_url ${BOS_INSTALL_ISO_URL} ${TEMP_FILE}
        ISO_DEVICE=$(losetup --show -f $TEMP_FILE)
        rm -f $TEMP_FILE
    fi

    if [ -z "${ISO_DEVICE}" ]; then
        echo "#### There is no bos ISO device"
        return 1
    fi
}

setup_style()
{
    if [ "$BOS_INSTALL_FORCE_EFI" = "true" ] || [ -e /sys/firmware/efi ]; then
        PARTTABLE=gpt
        BOOTFLAG=esp
        if [ ! -e /sys/firmware/efi ]; then
            echo WARNING: installing EFI on to a system that does not support EFI
        fi
    else
        PARTTABLE=msdos
        BOOTFLAG=boot
    fi
}

validate_progs()
{
    for i in $PROGS; do
        if [ ! -x "$(which $i)" ]; then
            MISSING="${MISSING} $i"
        fi
    done

    if [ -n "${MISSING}" ]; then
        echo "The following required programs are missing for installation: ${MISSING}"
        exit 1
    fi
}

validate_device()
{
    DEVICE=$BOS_INSTALL_DEVICE
    if [ ! -b ${DEVICE} ]; then
        echo "You should use an available device. Device ${DEVICE} does not exist."
        exit 1
    fi
}

create_opt()
{
    mkdir -p "${TARGET}/bos/data/opt"
}

while [ "$#" -gt 0 ]; do
    case $1 in
        --no-format)
            BOS_INSTALL_NO_FORMAT=true
            ;;
        --force-efi)
            BOS_INSTALL_FORCE_EFI=true
            ;;
        --poweroff)
            BOS_INSTALL_POWER_OFF=true
            ;;
        --takeover)
            BOS_INSTALL_TAKE_OVER=true
            ;;
        --debug)
            set -x
            BOS_INSTALL_DEBUG=true
            ;;
        --config)
            shift 1
            BOS_INSTALL_CONFIG_URL=$1
            ;;
        --tty)
            shift 1
            BOS_INSTALL_TTY=$1
            ;;
        -h)
            usage
            ;;
        --help)
            usage
            ;;
        *)
            if [ "$#" -gt 2 ]; then
                usage
            fi
            INTERACTIVE=true
            BOS_INSTALL_DEVICE=$1
            BOS_INSTALL_ISO_URL=$2
            break
            ;;
    esac
    shift 1
done

if [ -e /etc/environment ]; then
    source /etc/environment
fi

if [ -e /etc/os-release ]; then
    source /etc/os-release

    if [ -z "$BOS_INSTALL_ISO_URL" ]; then
        BOS_INSTALL_ISO_URL=${ISO_URL}
    fi
fi

if [ -z "$BOS_INSTALL_DEVICE" ]; then
    usage
fi

validate_progs
validate_device

trap cleanup exit

get_iso
setup_style
do_format
do_mount
do_copy
install_grub
create_opt

if [ -n "$INTERACTIVE" ]; then
    exit 0
fi

if [ "$BOS_INSTALL_POWER_OFF" = true ] || grep -q 'bos.install.power_off=true' /proc/cmdline; then
    poweroff -f
else
    echo " * Rebooting system in 5 seconds (CTRL+C to cancel)"
    sleep 5
    reboot -f
fi