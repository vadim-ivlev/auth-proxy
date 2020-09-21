#!/bin/bash

RSA_PRIVATE_KEY=$1

which ssh-agent || ( apt-get update -y && apt-get install openssh-client -y )
eval $(ssh-agent -s)
ssh-add <(echo "$RSA_PRIVATE_KEY" | base64 --decode)
mkdir -p /root/.ssh
echo -e "Host *\n\tStrictHostKeyChecking no\n\n" >~/.ssh/config
echo "$RSA_PRIVATE_KEY" | base64 --decode >/root/.ssh/id_rsa
chmod 600 /root/.ssh/id_rsa && chmod 700 /root/.ssh
