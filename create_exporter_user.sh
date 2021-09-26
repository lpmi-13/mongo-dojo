echo "sleeping for 10 seconds..."
sleep 10
mongo localhost:27017 <<EOF
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

echo "all done!"
