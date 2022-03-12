#! /bin/bash

vagrant up

docker build -t data-insert data-insert/

docker run -it --rm data-insert
