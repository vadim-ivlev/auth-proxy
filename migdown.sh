#!/bin/bash
# Удаляет объекты базы данных
migrate -source=file://migrations/ -database postgres://root:root@localhost:5432/onlinebc?sslmode=disable down
