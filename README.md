# Mongo dojo


There are going to be two different configurations for running these exercises:

1. Running a replicaset inside of VMs

This is a basic setup to reproduce a simple mongo replicaset using VMs locally with virtualbox and vagrant. I chose virtualbox because I figured that would be better cross-platform than something like parallels, and also because one of the things we'd like to practice is killing and restarting the `mongod` process within the VM, which is less straightforward in containers.

2. Running a replicaset inside containers

This is all containerized, and it's used for the tasks that don't actually involve stopping the mongod process, which normally would kill the container. The dashboard config is a bit different, since we can just use the container network DNS to address each of the processes instead of IP address, but everything else is basically the same.

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
(this just starts up all the containers and doesn't insert any data yet)

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


- Creating an index



- Load data and start with an unoptimized query. Run explain to see that it's not performant and fix it. Re-run explain to prove the fix worked.

- Running `explain()` to see why a query might be running slowly (probably lack of an index)

- Find and stop a long-running query

- Firing exactly enough traffic to the primary to make it stop responding to read requests, but still respond to heartbeats

- Attaching a node service that reads from the primary (without secondaryPreferred), and then see what happens when the primary steps down (it should break)

- Fire a number of different types of queries into mongo and see what the graphs look like: skip param with a high number (1000+), gt/ls combined in the same query maybe?

