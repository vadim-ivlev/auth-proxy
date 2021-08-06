#!/bin/bash

# Если под Windows, добавляем команду sudo
# if [[ "$OSTYPE" == "msys" ]]; then alias sudo=""; fi


echo 'гасим бд'
docker-compose -f docker-compose-dev.yml down -v


