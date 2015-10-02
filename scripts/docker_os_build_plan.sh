#!/bin/bash
echo "This script will pre-configure the server to run docker"
DOCKER_USER_INPUT=$1
DOCKER_PROXY_INPUT=$2
DOCKER_PUBKEY_INPUT=$3
if [ -z "${DOCKER_PUBKEY_INPUT}" ]; then
  echo "ERROR : this script requires a public key for docker user!"
  echo "USAGE: $0 <docker user> '<public key>'"
  exit 1
fi

DOCKER_USER=${DOCKER_USER_INPUT:-"docker"}
DOCKER_PROXY=${DOCKER_PROXY_INPUT}
DOCKER_PUBKEY=${DOCKER_PUBKEY_INPUT}

# boot the external interface, replace this to another interface dependening on your hardware
ifup eno50

# optionally set first boot proxy server configuration
export http_proxy="http://proxy:8080"
export https_proxy="http://proxy:8080"
export HTTP_PROXY="http://proxy:8080"
export HTTPS_PROXY="http://proxy:8080"

# optionally set some persistent proxy server configuration
cat >> "/root/.bash_profile" << EOF
${DOCKER_PROXY}
EOF

# create a service account
useradd "${DOCKER_USER}" -d "/home/${DOCKER_USER}"

# setup .ssh folderls -ak
if [ ! -d "/home/${DOCKER_USER}/.ssh" ]; then
  mkdir -p "/home/${DOCKER_USER}/.ssh"
  chmod 700 "/home/${DOCKER_USER}/.ssh"
  chown "${DOCKER_USER}:${DOCKER_USER}" "/home/${DOCKER_USER}/.ssh"
fi
if [ ! -f "/home/${DOCKER_USER}/.ssh/authorized_keys" ] ; then
  touch "/home/${DOCKER_USER}/.ssh/authorized_keys"
  chmod 600 "/home/${DOCKER_USER}/.ssh/authorized_keys"
  chown "${DOCKER_USER}:${DOCKER_USER}" "/home/${DOCKER_USER}/.ssh/authorized_keys"
fi
cat >> "/home/${DOCKER_USER}/.ssh/authorized_keys" << EOF
${DOCKER_PUBKEY}
EOF

# modify /home/{user}/.bash_profile to set a persistent proxy
cat >> "/home/${DOCKER_USER}/.bash_profile" << EOF
${DOCKER_PROXY}
EOF

# give sudoers access
cat >> "/etc/sudoers.d/90-${DOCKER_USER}" << SUDOERS_EOF
# User rules for icsp docker user
${DOCKER_USER} ALL=(ALL) NOPASSWD:ALL
SUDOERS_EOF

# modify primary nic eno50 to start on boot
sed -i 's/ONBOOT=no/ONBOOT=yes/g' /etc/sysconfig/network-scripts/ifcfg-eno50
sed -i 's/localhost.localdomain/egslcloud-scs79.fc.hpe.com/g' /etc/hostname
shutdown -r now