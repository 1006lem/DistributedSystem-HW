#!/bin/bash

# Check operating system
if [[ $(uname -s) != "Linux" ]]; then
    echo "This script is for Linux."
    exit 1
fi


cd $HOME
curl -LO https://go.dev/dl/go1.20.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
cat ~/.profile | grep /usr/local/go/bin || echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
cat ~/.profile | grep GOPATH= || echo 'export GOPATH=$(go env GOPATH)' >> ~/.profile
cat ~/.profile | grep GOPATH/bin || echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.profile
source ~/.profile

echo $GOPATH
go version
