#!/bin/bash

# echo "building ..."
# export GO111MODULE=on

# env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a  .

# линкуем статически под линукс
# env CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' . 


echo "cleaning  deploy/ directory "

# clean deploy/ directory
rm -rf deploy/auth-proxy
# rm -rf deploy/configs_example
rm -rf deploy/migrations
rm -rf deploy/templates
rm -rf deploy/etc
rm -rf deploy/certificates

# careful with configs, everything excluding db.yaml
rm -f deploy/configs/mail.yaml
rm -f deploy/configs/mail-templates.yaml.yaml
rm -f deploy/configs/app.yaml
rm -f deploy/configs/oauth2.yaml
rm -f deploy/configs/signature.yaml


echo "copying stuff ..."
# copy files to deploy/
cp    auth-proxy        deploy/auth-proxy
# cp -r configs           deploy/configs_example
cp -r migrations        deploy/migrations
cp -r templates         deploy/templates
cp -r etc               deploy/etc
cp -r certificates      deploy/certificates


# careful with configs, everything excluding db.yaml
mkdir -p deploy/configs
cp -f configs/gl/mail.yaml  deploy/configs/mail.yaml
cp -f configs/gl/mail-templates.yaml  deploy/configs/mail-templates.yaml
cp -f configs/gl/app.yaml  deploy/configs/app.yaml
cp -f configs/gl/oauth2.yaml  deploy/configs/oauth2.yaml
cp -f configs/gl/signature.yaml  deploy/configs/signature.yaml
