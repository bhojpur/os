
variable "box_description" {
  type    = string
  default = "Bhojpur OS is a Linux distribution designed to remove as much OS maintenance as possible in a Kubernetes cluster"
}

variable "box_version" {
  type    = string
  default = "v0.20.7-dcp1r0"
}

variable "iso_checksum" {
  type    = string
  default = "85a560585bc5520a793365d70e6ce984f3fb2ce5a43b31f0f7833dc347487e69"
}

variable "iso_url" {
  type    = string
  default = "https://github.com/bhojpur/os/releases/download/v0.20.7-dcp1r0/bos-amd64.iso"
}

variable "password" {
  type    = string
  default = "bhojpur"
}

source "virtualbox-iso" "bos" {
  boot_command = [
    "bhojpur", "<enter>",
    "sudo opsutl install", "<enter>",
    "1", "<enter>",
    "y", "<enter>",
    "http://{{ .HTTPIP }}:{{ .HTTPPort }}/config.yml", "<enter>",
    "y", "<enter>",
  ]
  boot_wait            = "40s"
  disk_size            = "8000"
  export_opts          = ["--manifest", "--vsys", "0", "--description", "${var.box_description}", "--version", "${var.box_version}"]
  format               = "ova"
  guest_os_type        = "Linux_64"
  http_directory       = "."
  iso_checksum         = "sha256:${var.iso_checksum}"
  iso_url              = "${var.iso_url}"
  post_shutdown_delay  = "10s"
  shutdown_command     = "sudo poweroff"
  ssh_keypair_name     = ""
  ssh_private_key_file = "packer_rsa"
  ssh_timeout          = "1000s"
  ssh_username         = "bhojpur"
}

build {
  sources = ["source.virtualbox-iso.bos"]

  post-processor "vagrant" {
    output = "bos_{{.Provider}}.box"
  }
}