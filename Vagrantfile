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
sudo /etc/init.d/network restart
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
  end
end
