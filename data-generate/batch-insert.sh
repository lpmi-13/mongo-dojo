#!/bin/bash

# this argument should either be "localhost" if you're using containers, or "192.168.42.102" for VMs
HOSTNAME=$1

if [[ -z "$HOSTNAME" ]]; then
  echo "please pass in a hostname as the first argument, like\n $ ./batch-import.sh localhost"
  exit 1
fi

for record in $(find . -type f -name "record-*"); do
    echo entering $record ...
    mongoimport --uri="mongodb://$HOSTNAME:27017/userData" --jsonArray --collection=reviews $record
done;
