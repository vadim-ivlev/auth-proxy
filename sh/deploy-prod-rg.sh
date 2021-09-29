#!/bin/bash

echo "deploying on production rg ..."
git push origin master:production
echo "Watch pipeline -> https://git.rgwork.ru/masterback/auth-proxy/pipelines"

