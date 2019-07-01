#!/bin/bash

# линкуем статически под линукс
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a  .

# env CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' . 