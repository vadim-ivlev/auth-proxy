#!/bin/bash

echo "deploying ..."
git push origin master:production
echo "Watch pipeline -> https://git.rgwork.ru/masterback/auth-proxy/pipelines"

