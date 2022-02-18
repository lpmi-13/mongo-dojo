#!/bin/sh

docker build -t data-insert-containers data-insert-containers/

# we need to run this on the host network, since this isn't in the compose network, which also means we need to manually add the hosts mapping
docker run -it --rm --network host --add-host mongo1:127.0.0.1 data-insert-containers
