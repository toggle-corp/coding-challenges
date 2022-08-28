#!/bin/bash

cat /keys/id_rsa.pub >> /home/runner/.ssh/authorized_keys
chown runner:runner -R /home/runner/.ssh
chmod 600 /home/runner/.ssh/authorized_keys

/usr/sbin/sshd -D
