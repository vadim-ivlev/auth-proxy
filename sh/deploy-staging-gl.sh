#!/bin/bash

echo "deploying on staging gl ..."
git push origin master:staging-gl
echo "Watch pipeline -> https://git.rgwork.ru/masterback/auth-proxy/pipelines"

