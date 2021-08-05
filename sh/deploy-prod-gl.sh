#!/bin/bash

echo "deploying on production gl ..."
git push origin master:production-gl
echo "Watch pipeline -> https://git.rgwork.ru/masterback/auth-proxy/pipelines"
