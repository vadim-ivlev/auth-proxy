#!/bin/bash

echo 'поднимаем докер'
docker-compose -f docker-compose-public.yml up -d && docker logs -f auth-proxy-public 
