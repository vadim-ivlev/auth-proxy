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
rm -rf deploy/configs_example
rm -rf deploy/migrations
rm -rf deploy/templates
rm -rf deploy/testapps

# careful with configs/
rm -f deploy/configs/mail.yaml
rm -f deploy/configs/sqlite.yaml



# copy files to deploy/
cp    auth-proxy    deploy/auth-proxy
cp -r configs       deploy/configs_example
cp -r migrations    deploy/migrations
cp -r templates     deploy/templates
cp -r testapps      deploy/testapps

mv deploy/node_modules      deploy/nodemodules


# careful with configs/
mkdir -p deploy/configs
cp -f configs/mail.yaml  deploy/configs/mail.yaml
cp -f configs/sqlite.yaml  deploy/configs/sqlite.yaml
