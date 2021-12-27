#!/bin/bash

# install docker
sudo yum update -y
sudo yum install -y docker
# install git make
sudo yum install -y git make wget

# install docker-compose
sudo wget https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m) 
sudo mv docker-compose-$(uname -s)-$(uname -m) /usr/bin/docker-compose
sudo chmod -v +x /usr/bin/docker-compose
sudo chmod -v +x /usr/bin/docker-compose

# enable docker
sudo systemctl enable docker.service
sudo systemctl start docker.service

# get moodle docker-compose
git clone https://github.com/bitnami/bitnami-docker-moodle.git

echo "cd bitnami-docker-moodle/3/debian-10"
cd bitnami-docker-moodle/3/debian-10

# change settings for japanese
sed -i -e "29i RUN echo 'ja_JP.UTF-8 UTF-8' >> /etc/locale.gen && locale-gen" ./Dockerfile

# echo "sudo docker-compose up -d"
sudo docker-compose build
# sudo docker-compose up -d
