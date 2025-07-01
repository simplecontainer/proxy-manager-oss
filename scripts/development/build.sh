#!/bin/bash
cd "$(dirname "$0")" || exit 1
cd ../../

echo "Doing work in directory $PWD"

BASE_DIR="$PWD"

cd "$BASE_DIR" || exit 1

echo "***********************************************"
echo "$(pwd)"
echo "***********************************************"

CGO_ENABLED=0 go build -ldflags '-s -w' || exit 1

mkdir $BASE_DIR/proxy-manager-oss-linux-amd64 || echo "$BASE_DIR/proxy-manager-oss already created"
cp -f proxy-manager-oss $BASE_DIR/proxy-manager-oss-linux-amd64/proxy-manager-oss

cd "$BASE_DIR" || exit 1