#!/bin/bash

vagrant up

cd data-generate && docker-compose up --build && ./batch-insert.sh 192.168.42.102
