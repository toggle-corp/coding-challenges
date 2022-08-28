#!/bin/bash

ssh-keyscan runnerenv > ~/.ssh/known_hosts
/go/app/scripts/wait-for-it.sh $POSTGRES_HOSTNAME:$POSTGRES_PORT && \
    /go/app/coding-challenge
