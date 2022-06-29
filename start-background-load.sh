#!/bin/bash

cd data-query-background && docker-compose up -d --build

cd ../data-insert-background && docker-compose up -d --build
