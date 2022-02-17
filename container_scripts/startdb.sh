#!/bin/bash

docker-compose up -d

sleep 5

docker exec mongo1 /scripts/rs-init.sh

sleep 15

docker exec mongo1 /scripts/create_exporter_user.sh
