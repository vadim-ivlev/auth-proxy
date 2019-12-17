#!/bin/bash

echo "building ..."
export GO111MODULE=on

# env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a  .

# линкуем статически под линукс
# env CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' . 

# CGO_ENABLED=1 нужно для драйвера SQLite
env CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' .  || exit 1


echo "cleaning  deploy/ directory "

# clean deploy/ directory
rm -rf deploy/auth-proxy
rm -rf deploy/configs_example
rm -rf deploy/migrations
rm -rf deploy/templates
rm -rf deploy/node-apps
rm -rf deploy/etc
rm -rf deploy/certificates

# careful with configs, everything excluding db.yaml
rm -f deploy/configs/mail.yaml
rm -f deploy/configs/mail-templates.yaml.yaml
rm -f deploy/configs/sqlite.yaml
rm -f deploy/configs/app.yaml
rm -f deploy/configs/oauth2.yaml
rm -f deploy/configs/signature.yaml


echo "copying stuff ..."
# copy files to deploy/
cp    auth-proxy        deploy/auth-proxy
cp -r configs           deploy/configs_example
cp -r migrations        deploy/migrations
cp -r templates         deploy/templates
cp -r node-apps         deploy/node-apps
cp -r etc               deploy/etc
cp -r certificates      deploy/certificates

#mv deploy/node-apps/node_modules      deploy/node-apps/nodemodules


# careful with configs, everything excluding db.yaml
mkdir -p deploy/configs
cp -f configs/mail.yaml  deploy/configs/mail.yaml
cp -f configs/mail-templates.yaml  deploy/configs/mail-templates.yaml
cp -f configs/sqlite.yaml  deploy/configs/sqlite.yaml
cp -f configs/app.yaml  deploy/configs/app.yaml
cp -f configs/oauth2.yaml  deploy/configs/oauth2.yaml
cp -f configs/signature.yaml  deploy/configs/signature.yaml
