#!/bin/bash

HOME="/home/ubuntu"
WORKSPACE="${HOME}/go/src/github.com/minchao/smsender"
CONFIG="${WORKSPACE}/config/config.yml"
MYSQL_USER=${MYSQL_USER:-smsender_user}
MYSQL_PASSWORD=${MYSQL_PASSWORD:-smsender_password}
MYSQL_HOST=${MYSQL_HOST:-127.0.0.1}
MYSQL_PORT=${MYSQL_PORT:-3306}
MYSQL_DATABASE=${MYSQL_DATABASE:-smsender}

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
ln -s /vagrant "${WORKSPACE}"

echo "Setting up development environment"

echo 'export PATH=$PATH:/usr/local/go/bin:~/go/bin' >> "${HOME}/.profile"
echo '# Automatically chdir to workspace upon vagrant ssh' >> "${HOME}/.profile"
echo "cd ${WORKSPACE}" >> "${HOME}/.profile"
source "${HOME}/.profile"

cd "${WORKSPACE}"

echo "Build app"

make build

echo "Configure config.yml"
if [ ! -f ${CONFIG} ]
then
  cp ./config/config.default.yml ${CONFIG}
  sed -Ei "s/MYSQL_USER/$MYSQL_USER/" ${CONFIG}
  sed -Ei "s/MYSQL_PASSWORD/$MYSQL_PASSWORD/" ${CONFIG}
  sed -Ei "s/MYSQL_HOST/$MYSQL_HOST/" ${CONFIG}
  sed -Ei "s/MYSQL_PORT/$MYSQL_PORT/" ${CONFIG}
  sed -Ei "s/MYSQL_DATABASE/$MYSQL_DATABASE/" ${CONFIG}
  echo OK
fi

echo "Starting docker-compose-dev.yml"

sudo docker-compose -f docker-compose-dev.yml up
