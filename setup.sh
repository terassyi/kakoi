#!/bin/bash

mkdir -p /etc/kakoi/templates
cp -r ./templates/* /etc/kakoi/templates

go build .
mv kakoi /usr/local/bin/kakoi