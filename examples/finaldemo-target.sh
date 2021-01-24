#!/bin/bash

sudo yum update -y
sudo yum install -y git nmap

git clone https://github.com/terassyi/demokey.git
cp demokey/* ~/.ssh/
chmod 0600 ~/.ssh/id_rsa
chmod 0600 ~/.ssh/id_rsa.pub
cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys

echo "this is demo flag!!!!" > flag.txt