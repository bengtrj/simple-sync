#!/usr/bin/env bash

set +eu

sudo apt-get update
sudo apt-get -y upgrade
wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz
sudo tar -xvf go1.13.3.linux-amd64.tar.gz
sudo mv go /usr/local

export GOROOT=/usr/local/go
export PATH=$GOROOT/bin:$PATH

./make.sh