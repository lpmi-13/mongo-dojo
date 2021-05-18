# -*- mode: ruby -*-
# vi: set ft=ruby :

$mongoConfigScript = <<-"SCRIPT"
  YUM_REPO_CONFIG_PATH="/etc/yum.repos.d/mongodb.repo"

  tee $YUM_REPO_CONFIG_PATH <<-"EOF"

[mongodb-org-3.6]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/3.6/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-3.6.asc
EOF

sudo yum install -y mongodb-org

MONGOD_CONF_FILE="/etc/mongod.conf"

tee $MONGOD_CONF_FILE <<-"EOF"
# mongod.conf
systemLog:
  destination: file
  logAppend: true
  path: /var/log/mongodb/mongod.log
# Where and how to store data.
storage:
  dbPath: /var/lib/mongo
  journal:
    enabled: true
# how the process runs
processManagement:
  fork: true  # fork and run in background
  pidFilePath: /var/run/mongodb/mongod.pid  # location of pidfile
  timeZoneInfo: /usr/share/zoneinfo
# network interfaces
net:
  port: 27017
  bindIp: 0.0.0.0  # Enter 0.0.0.0,:: to bind to all IPv4 and IPv6 addresses or, alternatively, use the net.bindIpAll setting.
replication:
   oplogSizeMB: 50
   replSetName: dojo
EOF

sudo systemctl restart mongod
sudo systemctl enable mongod

# I thought there would be a cleaner way to do this, but this works so keeping it for now
sudo iptables -A IN_public_allow -p tcp -m tcp --dport 27017 -m conntrack --ctstate NEW,UNTRACKED -j ACCEPT
sudo iptables -A IN_public_allow -p tcp -m tcp --dport 9216 -m conntrack --ctstate NEW,UNTRACKED -j ACCEPT
sudo /etc/init.d/network restart



wget https://github.com/percona/mongodb_exporter/releases/download/v0.7.1/mongodb_exporter-0.7.1.linux-amd64.tar.gz
tar xvzf mongodb_exporter-0.7.1.linux-amd64.tar.gz
sudo mv mongodb_exporter /usr/local/bin/
sudo useradd -rs /bin/false prometheus

export MONGODB_URI=mongodb://mongodb_exporter:password@localhost:27017

MONGODB_EXPORTER_PATH="/lib/systemd/system/mongodb_exporter.service"

tee $MONGODB_EXPORTER_PATH <<-"EOF"
[Unit]
Description=MongoDB Exporter
User=prometheus

[Service]
Type=simple
Restart=always
ExecStart=/usr/local/bin/mongodb_exporter

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl start mongodb_exporter.service


SCRIPT

$observerConfigScript = <<-"SCRIPT"

wget https://github.com/prometheus/prometheus/releases/download/v2.9.2/prometheus-2.9.2.linux-amd64.tar.gz
tar xvzf prometheus-2.9.2.linux-amd64.tar.gz


PROMETHEUS_CONFIG_PATH="/etc/prometheus/prometheus.yml"

sudo mkdir /etc/prometheus
sudo mkdir /var/lib/prometheus

tee $PROMETHEUS_CONFIG_PATH <<-"EOF"
# A scrape configuration containing exactly one endpoint to scrape:
global:
  scrape_interval: 1s

scrape_configs:
  - job_name: 'mongo_repl'
    static_configs:
      - targets: ['192.168.42.100:9216', '192.168.42.101:9216', '192.168.42.102:9216']
EOF

sudo useradd -rs /bin/false prometheus

sudo chown prometheus:prometheus /etc/prometheus
sudo chown prometheus:prometheus /var/lib/prometheus

sudo cp prometheus-2.9.2.linux-amd64/{prometheus,promtool} /usr/local/bin/
sudo cp -r prometheus-2.9.2.linux-amd64/consoles /etc/prometheus/consoles
sudo cp -r prometheus-2.9.2.linux-amd64/console_libraries /etc/prometheus/console_libraries

sudo chown -R prometheus:prometheus /etc/prometheus

rm -rf prometheus-2.9.2.linux-amd64*

# let the host talk to prometheus
sudo iptables -A IN_public_allow -p tcp -m tcp --dport 9090 -m conntrack --ctstate NEW,UNTRACKED -j ACCEPT
sudo iptables -A IN_public_allow -p tcp -m tcp --dport 3000 -m conntrack --ctstate NEW,UNTRACKED -j ACCEPT
sudo /etc/init.d/network restart

PROMETHEUS_SERVICE="/lib/systemd/system/prometheus.service"
sudo tee $PROMETHEUS_SERVICE <<-"EOF"
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
  --config.file /etc/prometheus/prometheus.yml \
  --storage.tsdb.path /var/lib/prometheus/ \
  --web.console.templates=/etc/prometheus/consoles/ \
  --web.console.libraries=/etc/prometheus/console_libraries

[Install]
WantedBy=multi-user.target
EOF


  YUM_REPO_CONFIG_PATH="/etc/yum.repos.d/grafana.repo"

  tee $YUM_REPO_CONFIG_PATH <<-"EOF"

[grafana]
name=grafana
baseurl=https://packages.grafana.com/oss/rpm
repo_gpgcheck=1
enabled=1
gpgcheck=1
gpgkey=https://packages.grafana.com/gpg.key
sslverify=1
sslcacert=/etc/pki/tls/certs/ca-bundle.crt
EOF

sudo yum install grafana -y

GRAFANA_DATASOURCE="/etc/grafana/provisioning/datasources/default.yml"

tee $GRAFANA_DATASOURCE <<-"EOF"
apiVersion: 1

datasources:
  - name: Mongo Exporter
    type: mongo_metrics
EOF

GRAFANA_DEFAULT_DASHBOARD="/etc/grafana/provisioning/dashboards/default.yml"

tee $GRAFANA_DEFAULT_DASHBOARD <<-"EOF"
apiVersion: 1

providers:
  - name: Mongo_Exporter    # A uniquely identifiable name for the provider
    folder: Mongo_Stuff # The folder where to place the dashboards
    type: file
    options:
      path: /var/lib/grafana/dashboards

EOF

# hack because vagrant ssh user can't scp to /var
sudo mkdir -p /var/lib/grafana/dashboards
sudo cp /tmp/mongo_replicaset_summary.json /var/lib/grafana/dashboards/

sudo systemctl daemon-reload
sudo systemctl start prometheus
sudo systemctl enable prometheus
sudo systemctl start grafana-server
sudo systemctl enable grafana-server.service

sudo grafana-cli plugins install natel-discrete-panel
sudo grafana-cli plugins install digiapulssi-breadcrumb-panel
sudo grafana-cli plugins install yesoreyeram-boomtable-panel

sudo systemctl restart grafana-server

SCRIPT

Vagrant.configure("2") do |config|
  config.vm.provider "virtualbox" do |v|
    v.memory = 1024
    v.cpus = 1
  end
  config.vm.box = "generic/centos7"

  config.vm.define :mongo1 do |mongo1|
    mongo1.vm.network :private_network, ip: "192.168.42.100"
    mongo1.vm.hostname = "mongo1"
    mongo1.vm.provision "shell", inline: $mongoConfigScript
  end

  config.vm.define :mongo2 do |mongo2|
    mongo2.vm.network :private_network, ip: "192.168.42.101"
    mongo2.vm.hostname = "mongo2"
    mongo2.vm.provision "shell", inline: $mongoConfigScript
  end

  config.vm.define :mongo3 do |mongo3|
    mongo3.vm.network :private_network, ip: "192.168.42.102"
    mongo3.vm.hostname = "mongo3"
    mongo3.vm.provision "shell", inline: $mongoConfigScript
    mongo3.vm.provision "shell", path: "mongo_rs_config.sh"
    mongo3.vm.provision "shell", path: "create_exporter_user.sh"
  end

  config.vm.define :observer do |observer|
    observer.vm.network :private_network, ip: "192.168.42.200"
    observer.vm.network "forwarded_port", guest: 3000, host: 3000
    observer.vm.hostname = "observer"
    observer.vm.provision "file", source: "mongo_replicaset_summary.json", destination: "/tmp/mongo_replicaset_summary.json"
    observer.vm.provision "shell", inline: $observerConfigScript
  end
end
