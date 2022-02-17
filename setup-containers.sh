#!/bin/sh

docker build -t data-insert-containers data-insert-containers/

# we need to run this on the host network, since this isn't in the compose network
docker run -it --rm --network host data-insert-containers
