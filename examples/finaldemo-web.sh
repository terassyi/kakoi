#!/bin/bash

echo "yum update -y"
sudo yum update -y
echo "install git nmap"
sudo yum install -y git nmap

sudo yum search pytohn3
#echo "install python3"
#sudo yum installl -y python3
echo "install python3-pip"
sudo yum install -y python3-pip

git clone https://github.com/terassyi/vlunapp.git

git clone https://github.com/terassyi/demokey.git
cp demokey/* ~/.ssh/
chmod 0600 ~/.ssh/id_rsa
chmod 0600 ~/.ssh/id_rsa.pub
cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys

cd vlunapp

echo "install dependency"
sudo pip3 install -r requirements.txt

sudo cp /usr/local/bin/gunicorn /usr/bin/

echo "run app"
sudo gunicorn app:app --workers 2 --timeout 2 -b 0.0.0.0:80 -D