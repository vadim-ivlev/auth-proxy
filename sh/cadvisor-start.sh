#!/bin/bash


# github https://github.com/google/cadvisor

# use the latest release version from https://github.com/google/cadvisor/releases or "latest"
VERSION=v0.47.0 
sudo docker run --rm --name cadvisor \
  --volume=/:/rootfs:ro \
  --volume=/var/run:/var/run:ro \
  --volume=/sys:/sys:ro \
  --volume=/var/lib/docker/:/var/lib/docker:ro \
  --volume=/dev/disk/:/dev/disk:ro \
  --publish=8080:8080 \
  --detach=true \
  --name=cadvisor \
  --privileged \
  --device=/dev/kmsg \
  gcr.io/cadvisor/cadvisor:$VERSION

echo "open http://localhost:8080/"