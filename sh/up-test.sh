#!/bin/bash

echo 'поднимаем докер для front-end разработчиков'
docker-compose -f docker-compose-test.yml up -d && docker logs -f auth-proxy-front
