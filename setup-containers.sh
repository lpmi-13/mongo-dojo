#! /bin/bash

./container-scripts/startdb.sh

cd data-generate && docker-compose up --build && ./batch-insert.sh localhost