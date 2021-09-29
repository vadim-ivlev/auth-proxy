#!/bin/bash

echo "deploying on staging rg ..."
git push origin master:staging
echo "Watch pipeline -> https://git.rgwork.ru/masterback/auth-proxy/pipelines"

