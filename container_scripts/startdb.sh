#!/bin/bash

# clean out the mongo data locally first
rm -rf ~/mongors

docker-compose up -d

sleep 5

docker exec mongo1 /scripts/rs-init.sh

sleep 15

docker exec mongo1 /scripts/create_exporter_user.sh
