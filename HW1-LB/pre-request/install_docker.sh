#!/bin/bash

if [[ $(uname -s) != "Linux" ]]; then
    echo "This script is for Linux."
    exit 1
fi

if command -v docker &> /dev/null; then
    echo "Docker is already installed."
    exit 1

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install docker-ce
# configure daemon.json
sudo mkdir -p /etc/docker
cat <<EOF | sudo tee /etc/docker/daemon.json
{
    "exec-opts": ["native.cgroupdriver=systemd"],
    "log-driver": "json-file",
    "log-opts": {
        "max-size": "100m"
    },
    "storage-driver": "overlay2"
}
EOF

  sudo usermod -aG docker $USER
  sudo chmod 666 /var/run/docker.sock
fi

echo "Docker installed .."