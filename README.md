# Mongo dojo

This is a basic setup to reproduce a simple mongo replicaset using VMs locally with virtualbox and vagrant. I chose virtualbox because I figured that would be better cross-platform than something like parallels, and also because one of the things we'd like to practice is killing and restarting the `mongod` process within the VM, which is less straightforward in containers.

## Setting up

1) install Virtualbox
https://www.virtualbox.org/wiki/Downloads

2) install Vagrant
https://www.vagrantup.com/downloads

3) `$ vagrant box add generic/centos7 --provider virtualbox`

4) `$ vagrant plugin install vagrant-vbguest`

5) `$ bash setup.sh`
(this scripts insertion of mock data into the replicaset)

## Scenarios we want

### Basic housekeeping

- Upgrading the Mongo version
This will involve following the process in `steps.txt`.

- Restoring from backups (not exactly sure how backups will work in this setup, but we can probably work something out)

- look at "realtime" metrics inside the VM using "mongotop" and "mongostat"

### Working with indexes

- Creating an index

### Optimizing queries

- load data and start with an unoptimized query. Run explain to see that it's not performant and fix it. Re-run explain to prove the fix worked.

- Running `explain()` to see why a query might be running slowly (probably lack of an index)

- Find and stop a long-running query

### Seeing it break

- Firing exactly enough traffic to the primary to make it stop responding to read requests, but still respond to heartbeats

- Attaching a node service that reads from the primary (without secondaryPreferred), and then see what happens when the primary steps down (it should break)

- fix various things (probably a sub-set of the above) by restarting the `mongod` process

- fire a number of different types of queries into mongo and see what the graphs look like: skip param with a high number (1000+), gt/ls combined in the same query maybe?
