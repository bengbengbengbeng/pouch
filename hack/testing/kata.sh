#!/usr/bin/env bash

# integration::install_daishu installs qemu, guest kernel
# and image with vm needs
integration::install_daishu() {
  # change network mode
  sed -i 's/internetworking_model=\"enlightened"/internetworking_model=\"none\"/' /etc/kata-containers/configuration.toml || true

  # install daishu
  # install daishu only in physical machine
  if dmesg | grep -q "Hypervisor detected";then
      return
  fi

  if [ "$(rpm -q daishu)" != "" ];then
	  rpm -e "$(rpm -q daishu)" || true
  fi

  yum update -b test -y daishu-runtime || true

  local qemu kernel image
  qemu=$(grep -A 10 "hypervisor.qemu" /etc/kata-containers/configuration.toml | grep "^path" | awk -F\" '{print $2}')
  kernel=$(grep -A 10 "hypervisor.qemu" /etc/kata-containers/configuration.toml | grep "^kernel" | awk -F\" '{print $2}')
  image=$(grep -A 10 "hypervisor.qemu" /etc/kata-containers/configuration.toml | grep "^image" | awk -F\" '{print $2}')

  if [ -f "$qemu" ] && [ -f "$kernel" ] && [ -f "$image" ];then
      return
  fi

  echo "install daishu package"
  yum install -b test -y daishu-containers-image daishu-linux-container daishu-qemu
}
