#!/bin/bash

echo "deploying on production rg ..."
date=$(date '+%d%m%Y-%H%M')
git tag prod_$date
git push origin prod_$date
echo "Watch pipeline -> https://git.rgwork.ru/masterback/auth-proxy/pipelines"
