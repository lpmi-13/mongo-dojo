# Mongo dojo


There are going to be two different configurations for running these exercises:

1. Running a replicaset inside of VMs

This is a basic setup to reproduce a simple mongo replicaset using VMs locally with virtualbox and vagrant. I chose virtualbox because I figured that would be better cross-platform than something like parallels, and also because one of the things we'd like to practice is killing and restarting the `mongod` process within the VM, which is less straightforward in containers.

2. Running a replicaset inside containers

This is all containerized, and it's used for the tasks that don't actually involve stopping the mongod process, which normally would kill the container. The dashboard config is a bit different, since we can just use the container network DNS to address each of the processes instead of IP address, but everything else is basically the same.

## Generating the data

There are two ways to do this, one inserting the records as they're generated via Faker, and separating data creation from insertion into the mongo replicaset. Creating the json locally via `docker-compose up --build` in the `data-generate` directory, and then inserting via `./batch-insert.sh` takes about 10-12 minutes, whereas directly running the `./start_insert_swarm_containers.sh` takes almost 40 minutes.

## Setting up (for VMs)

1) install Virtualbox
https://www.virtualbox.org/wiki/Downloads

2) install Vagrant
https://www.vagrantup.com/downloads

3) `$ vagrant box add generic/centos7 --provider virtualbox`

4) `$ vagrant plugin install vagrant-vbguest`

5) `$ bash setup.sh`
(this scripts insertion of mock data into the replicaset)

### Tasks for the VM-based configuration

- Upgrading the Mongo version
This will involve following the process in `steps.txt`.

- Restoring from backups (not exactly sure how backups will work in this setup, but we can probably work something out)

We're going to use `vagrant snapshot save mongo3 backup` to snapshot the primary as a backup named "backup".

Then we're going to kill a machine (one of the secondaries) and restore that from the backup.

## Setting up (for containers)

1) Install Docker
https://docs.docker.com/get-docker/

2) Install Docker compose
https://docs.docker.com/compose/install/

3) `$ bash startdb.sh`

4) `$ bash start_insert_swarm_containers.sh`

### Tasks for the container-based configuration

- Look at "realtime" metrics inside the containers using "mongotop" and "mongostat"

```
docker exec -it mongo1 /bin/sh
```

and then once inside the container shell, just run `mongotop` and you should see something like:

```
# mongotop
2022-02-15T16:00:12.924+0000	connected to: 127.0.0.1

                    ns    total    read    write    2022-02-15T16:00:13Z
        local.oplog.rs      2ms     2ms      0ms
     admin.system.keys      0ms     0ms      0ms
    admin.system.roles      0ms     0ms      0ms
    admin.system.users      0ms     0ms      0ms
  admin.system.version      0ms     0ms      0ms
config.system.sessions      0ms     0ms      0ms
   config.transactions      0ms     0ms      0ms
        config.version      0ms     0ms      0ms
              local.me      0ms     0ms      0ms
local.replset.election      0ms     0ms      0ms
```

you can do the exact same for `mongostat`...

```
# mongostat
insert query update delete getmore command dirty used flushes vsize   res qrw arw net_in net_out conn  set repl                time
    *0    *0     *0     *0       0     2|0  0.0% 0.1%       0 1.42G 78.0M 0|0 1|0   415b   63.4k   12 dojo  PRI Feb 15 16:01:03.439
    *0    *0     *0     *0       0     4|0  0.0% 0.1%       0 1.42G 78.0M 0|0 1|0   825b   64.7k   12 dojo  PRI Feb 15 16:01:04.442
    *0    *0     *0     *0       0     3|0  0.0% 0.1%       0 1.42G 78.0M 0|0 1|0   418b   63.9k   12 dojo  PRI Feb 15 16:01:05.437
    *0     2     *0     *0       1    24|0  0.0% 0.1%       0 1.42G 78.0M 0|0 1|0  8.21k    167k   17 dojo  PRI Feb 15 16:01:06.434
    *0    *0     *0     *0       0    10|0  0.0% 0.1%       0 1.42G 78.0M 0|0 1|0  2.73k    111k   12 dojo  PRI Feb 15 16:01:07.444
    *0    *0     *0     *0       0     3|0  0.0% 0.1%       0 1.42G 78.0M 0|0 1|0   418b   63.9k   12 dojo  PRI Feb 15 16:01:08.438
```

### Run a query on an unindexed field (we'll index it later)

The `reviewsubmitted` field isn't currently indexed, so any query operation using it is going to result in a COLLSCAN

connect to the primary:

```
mongo mongodb://localhost:27017
```

and once connected to the primary, you should be able to see the number of documents in the `reviews` collection in the `userData` database:

```
dojo:PRIMARY> show dbs
admin     0.000GB
config    0.000GB
local     0.116GB
userData  0.113GB
dojo:PRIMARY> use userData
switched to db userData
dojo:PRIMARY> show collections
reviews
```

now we can run explain for a query based on the unindexed field

```
dojo:PRIMARY> db.reviews.find({ reviewsubmitted: { $lt: "2015-01-01 00:00:00"}}).ex
plain()
```

which will show us that mongo had to scan the entire collection to retrieve the results:

```
{
        "queryPlanner" : {
                "plannerVersion" : 1,
                "namespace" : "userData.reviews",
                "indexFilterSet" : false,
                "parsedQuery" : {
                        "reviewsubmitted" : {
                                "$lt" : "2015-01-01 00:00:00"
                        }
                },
                "winningPlan" : {
                        "stage" : "COLLSCAN",
                        "filter" : {
                                "reviewsubmitted" : {
                                        "$lt" : "2015-01-01 00:00:00"
                                }
                        },
                        "direction" : "forward"
                },
                "rejectedPlans" : [ ]
        },
        "serverInfo" : {
                "host" : "57411b882bcb",
                "port" : 27017,
                "version" : "3.6.23",
                "gitVersion" : "d352e6a4764659e0d0350ce77279de3c1f243e5c"
        },
        "ok" : 1,
        "operationTime" : Timestamp(1645142923, 1),
        "$clusterTime" : {
                "clusterTime" : Timestamp(1645142923, 1),
                "signature" : {
                        "hash" : BinData(0,"AAAAAAAAAAAAAAAAAAAAAAAAAAA="),
                        "keyId" : NumberLong(0)
                }
        }
}
```

The important part of this is the `"stage" : "COLLSCAN",` line, showing that mongo had to scan the entire collection. We don't want to do this if we can avoid it.


### Creating an index

first, connect to the primary:

(run this from your local machine)
```
mongo mongodb://localhost:27017
```


create an index on the `reviewsubmitted` field like so:

```
db.reviews.createIndex({ "reviewsubmitted": -1 })
```

and the output should be something like

```
{
        "createdCollectionAutomatically" : false,
        "numIndexesBefore" : 1,
        "numIndexesAfter" : 2,
        "ok" : 1,
        "operationTime" : Timestamp(1645142260, 1),
        "$clusterTime" : {
                "clusterTime" : Timestamp(1645142260, 1),
                "signature" : {
                        "hash" : BinData(0,"AAAAAAAAAAAAAAAAAAAAAAAAAAA="),
                        "keyId" : NumberLong(0)
                }
        }
}
```

now we can run explain on the same query as above, and we'll see it did a much more efficient index scan.

```
dojo:PRIMARY> db.reviews.find({ reviewsubmitted: { $lt: "2015-01-01 00:00:00"}}).explain()
{
        "queryPlanner" : {
                "plannerVersion" : 1,
                "namespace" : "userData.reviews",
                "indexFilterSet" : false,
                "parsedQuery" : {
                        "reviewsubmitted" : {
                                "$lt" : "2015-01-01 00:00:00"
                        }
                },
                "winningPlan" : {
                        "stage" : "FETCH",
                        "inputStage" : {
                                "stage" : "IXSCAN",
                                "keyPattern" : {
                                        "reviewsubmitted" : -1
                                },
                                "indexName" : "reviewsubmitted_-1",
                                "isMultiKey" : false,
                                "multiKeyPaths" : {
                                        "reviewsubmitted" : [ ]
                                },
                                "isUnique" : false,
                                "isSparse" : false,
                                "isPartial" : false,
                                "indexVersion" : 2,
                                "direction" : "forward",
                                "indexBounds" : {
                                        "reviewsubmitted" : [
                                                "(\"2015-01-01 00:00:00\", \"\"]"
                                        ]
                                }
                        }
                },
                "rejectedPlans" : [ ]
        },
        "serverInfo" : {
                "host" : "57411b882bcb",
                "port" : 27017,
                "version" : "3.6.23",
                "gitVersion" : "d352e6a4764659e0d0350ce77279de3c1f243e5c"
        },
        "ok" : 1,
        "operationTime" : Timestamp(1645143090, 1),
        "$clusterTime" : {
                "clusterTime" : Timestamp(1645143090, 1),
                "signature" : {
                        "hash" : BinData(0,"AAAAAAAAAAAAAAAAAAAAAAAAAAA="),
                        "keyId" : NumberLong(0)
                }
        }
}
```

The important information here is the `"stage" : "IXSCAN",` line, showing us that mongo did an index scan, which is WAY more efficient than a full COLLSCAN.

### Create a rolling index across the replicaset

For this, we'll do something that's much more common in a production system, where we need to create an index, but not stop mongo while it's happening, and do it on one instance at a time.

> TL;DR - So we remove one secondary from the replicaset, add the index, then update the configuration on the primary so that it's hidden (means it won't get hit for reads, so stale data isn't an issue), and once it catches up with replication, make it visible again for the replicaset. Then we do the same on the other secondary, and finally step down the primary, and once it becomes a secondary, do the same on that new secondary.

First, connect to a secondary via:

```
vagrant ssh mongo1
```
and now we update the mongo config to take the instance out of the replicaset

```
sudo vim /etc/mongod.conf
```

and update the following configuration to take the instance out of the replicaset

```
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
  port: 27117 # <-- CHANGE THIS TO SOMETHING BESIDES 27017
  bindIp: 0.0.0.0  # Enter 0.0.0.0,:: to bind to all IPv4 and IPv6 addresses or, alternatively, use the net.bindIpAll setting.
#replication:   <-- COMMENT THIS LINE OUT
#   oplogSizeMB: 50  <-- COMMENT THIS LINE OUT
#   replSetName: dojo  <--COMMENT THIS LINE OUT
```

and now we can restart the mongod process like so

```
sudo systemctl mongod restart
```

and then start the mongo shell again via

```
mongo mongodb://localhost:27017
```

and we're ready to create the index. You should be able to run the following:

```
use userData
db.reviews.createIndex({ "BusinessId": 1})
```

and then close out of the shell and examine the mongo logs to watch the index being built (it'll be fast, since nothing's coming into this instance anymore).

```
sudo tail -f /var/log/mongodb/mongod.log
```

and you should see something like...

```
2022-03-11T16:45:17.409+0000 I INDEX    [conn1] build index on: userData.reviews properties: { v: 2, key: { BusinessId: 1.0 }, name: "BusinessId_1", ns: "userData.reviews" }
2022-03-11T16:45:17.409+0000 I INDEX    [conn1] 	 building index using bulk method; build may temporarily use up to 500 megabytes of RAM
2022-03-11T16:45:20.001+0000 I -        [conn1]   Index Build: 1777900/5000001 35%
2022-03-11T16:45:23.001+0000 I -        [conn1]   Index Build: 3584000/5000001 71%
2022-03-11T16:45:36.999+0000 I INDEX    [conn1] build index done.  scanned 5000001 total records. 19 secs
```

*On the Primary node* set the instance with the index to be hidden, so while it syncs it's not also serving any reads and the syncing can happen more quickly.

```
vagrant ssh mongo3
```

and enter the mongo shell

```
mongo mongodb://localhost:27017
```

then update the replicaset config so the secondary can sync in a hidden state.

```
conf = rs.config()
```

and you should be able to verify which member you're targeting with:

```
conf.members[0]
{
	"_id" : 0,
	"host" : "192.168.42.100:27017",
	"arbiterOnly" : false,
	"buildIndexes" : true,
	"hidden" : false,
	"priority" : 1,
	"tags" : {
		
	},
	"slaveDelay" : NumberLong(0),
	"votes" : 1
}
```

so now we update this node to be `hidden` and with `priority: 0` so it can't accidentally become the primary.

```
conf.members[0].hidden = true
conf.members[0].priority = 0
rs.reconfig(conf)
```

and you should see output like

```
{
	"ok" : 1,
	"operationTime" : Timestamp(1647018316, 1),
	"$clusterTime" : {
		"clusterTime" : Timestamp(1647018316, 1),
		"signature" : {
			"hash" : BinData(0,"AAAAAAAAAAAAAAAAAAAAAAAAAAA="),
			"keyId" : NumberLong(0)
		}
	}
}
```

Now we're ready to change the port on the first node back to the default so it can reconnect to the replicaset and sync, with its new shiny index in place.

> on the first node (mongo1)

```
sudo vim /etc/mongod.conf
```

and set the port back to `27017` and uncomment the following lines:

```
#replication:
#   oplogSizeMB: 50
#   replSetName: dojo
```

then restart the mongod process

```
sudo systemctl mongod restart
```

The secondary should now be reconnected to the replicaset and syncing. So go ahead and check the grafana dashboard at `localhost:3000` just to confirm that everything looks okay, and repeat the exact same process for the other secondary.

After finishing up this process on both of the secondaries, we're ready to step down the primary so it can become a secondary, and the whole process can be repeated on the new secondary (previous primary).

To do that, from the primary node, enter the mongo shell and run:

```
rs.stepDown(60)
```

and you should see in the terminal that it now says `dojo:SECONDARY>`, at which point, just run through the same steps from above.

After that's done, congratulations! You've now run a full rolling index on a MongoDB replicaset!!!




- Find and stop a long-running query

- Firing exactly enough traffic to the primary to make it stop responding to read requests, but still respond to heartbeats

- Attaching a node service that reads from the primary (without secondaryPreferred), and then see what happens when the primary steps down (it should break)

- Fire a number of different types of queries into mongo and see what the graphs look like: skip param with a high number (1000+), gt/ls combined in the same query maybe?


### Gotchas

- Grafana loses the connection to its datasource

If, at any point, you need to suspend the vagrant machines, you might need to reprovision the grafana/prometheus components to pick up the datasources again. You'll know this because the grafana dashboards will show "no data" for every panel.

```
vagrant destroy observer
```

and then

```
vagrant up --provision observer
```

- Prometheus can't scrape metrics from one of the nodes

If you see that the exporters are working locally on the node via `curl localhost:9100` and `curl localhost:9216`, but you can't hit those ports from outside the nodes, then you might just need to reset the iptables rules and restart the network via systemd on the node that's acting up.

```
sudo iptables -A IN_public_allow -p tcp -m tcp --dport 9216 -m conntrack --ctstate NEW,UNTRACKED -j ACCEPT
sudo iptables -A IN_public_allow -p tcp -m tcp --dport 9100 -m conntrack --ctstate NEW,UNTRACKED -j ACCEPT
sudo /etc/init.d/network restart
```
