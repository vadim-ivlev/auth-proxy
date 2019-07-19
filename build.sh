#!/bin/bash

echo "building ..."
export GO111MODULE=on

# линкуем статически под линукс
# env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a  .

# env CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' . 

env CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' .


echo "copying stuff ..."

# clean deploy/ directory
rm -rf deploy/auth-proxy
rm -rf deploy/configs
rm -rf deploy/migrations
rm -rf deploy/templates
rm -rf deploy/test_apps


# copy files to deploy/
cp    auth-proxy    deploy/auth-proxy
cp -r configs       deploy/configs
cp -r migrations    deploy/migrations
cp -r templates     deploy/templates
cp -r test_apps     deploy/test_apps
