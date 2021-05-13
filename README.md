# Mongo dojo

This is a basic setup to reproduce a simple mongo replicaset using VMs locally with virtualbox and vagrant. I chose virtualbox because I figured that would be better cross-platform than something like parallels.

## Setting up

1) install Virtualbox
https://www.virtualbox.org/wiki/Downloads

2) install Vagrant
https://www.vagrantup.com/downloads

3) `$ vagrant box add generic/centos7 --provider virtualbox`

4) `$ vagrant init generic/centos7`
(this is the step that sets up your skeleton `Vagrantfile`)

5) `$ vagrant plugin install vagrant-vbguest`

6) `$ vagrant up`

7) `$ vagrant ssh`

You're in your VM!

## Scenarios we want

- Upgrading the Mongo version

This will involve following the process in `steps.txt`.

- Manually stepping down the primary

- Running `explain()` to see why a query might be running slowly (probably lack of an index)

- Creating an index

- Restoring from backups (not exactly sure how backups will work in this setup, but we can probably work something out
