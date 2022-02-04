#!/bin/bash

echo 'гасим докер'
docker-compose -f docker-compose-public.yml down -v --remove-orphans 
