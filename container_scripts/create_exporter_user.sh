#!/bin/bash

mongo <<EOF
use admin;
rs.status();

# so the exporter can grab data and provide it for prometheus
db.createUser({
  user: "mongodb_exporter",
  pwd: "password",
  roles: [
      { role: "clusterMonitor", db: "admin" },
      { role: "read", db: "local" }
  ]
})
EOF
