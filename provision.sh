#!/bin/bash

HOME="/home/ubuntu"

sudo apt update
sudo apt upgrade -y
sudo apt install -y git curl make
curl -sSL https://get.docker.com/ | sudo sh
sudo usermod -aG docker ubuntu
sudo apt install -y docker-compose

echo "Installing Go"
if [ ! -e "/vagrant/go.tar.gz" ]; then
    curl https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz -o /vagrant/go.tar.gz
fi;
sudo tar -C /usr/local -xzf /vagrant/go.tar.gz

echo "Creating workspace"
mkdir -p "${HOME}/go/src/github.com/minchao"
ln -s /vagrant "${HOME}/go/src/github.com/minchao/smsender"

echo "Adding environment variable"
echo 'export PATH=$PATH:/usr/local/go/bin' >> "${HOME}/.profile"
echo '# Automatically chdir to workspace upon vagrant ssh' >> "${HOME}/.profile"
echo 'cd ~/go/src/github.com/minchao/smsender' >> "${HOME}/.profile"
source "${HOME}/.profile"
