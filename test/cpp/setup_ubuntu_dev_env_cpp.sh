#!/bin/bash

# This script is used to setting up the C++ dev env

echo "Installing required packages..."
sudo apt-get -y install build-essential autoconf libtool pkg-config
sudo apt-get -y install libgflags-dev libgtest-dev
sudo apt-get -y install clang libc++-dev
sudo apt-get -y install curl

echo "Building grpc from source..."
git clone -b $(curl -L https://grpc.io/release) https://github.com/grpc/grpc
cd grpc
git submodule update --init
sudo make install

echo "Installing protoc..."
cd ../grpc/third_party/protobuf
sudo make install

echo "***Installation complete***"
