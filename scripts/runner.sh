#!/bin/bash

lang=$1
fpath=$2

echo This message goes to stderr >&2
echo This also message goes to stderr >&2
exit 1
