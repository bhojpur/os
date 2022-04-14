ARG REPO
ARG TAG
FROM ${REPO}/bos-tar:${TAG} as tar

ARG REPO
ARG TAG
FROM ${REPO}/bos-iso:${TAG} as iso

ARG REPO
ARG TAG
FROM ${REPO}/bos-kernel:${TAG} as kernel

ARG REPO
ARG TAG
FROM ${REPO}/bos-base:${TAG}
ARG ARCH

COPY --from=kernel /output/vmlinuz /output/bos-vmlinuz-${ARCH}
COPY --from=kernel /output/initrd /output/bos-initrd-${ARCH}
COPY --from=kernel /output/kernel.squashfs /output/bos-kernel-${ARCH}.squashfs
COPY --from=kernel /output/version /output/bos-kernel-version-${ARCH}
COPY --from=iso /output/bos.iso /output/bos-${ARCH}.iso
COPY --from=tar /output/userspace.tar /output/bos-rootfs-${ARCH}.tar
RUN gzip /output/bos-rootfs-${ARCH}.tar