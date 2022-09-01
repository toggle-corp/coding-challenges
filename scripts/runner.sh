#!/bin/bash

lang=$1
submission_id=$2

# seconds
RUN_LIMIT=5
dirname=submission_$submission_id

cleanup () {
    rm -r /tmp/$dirname
}

# create directory
mkdir /tmp/$dirname

# get Extension from lang
ext=
runcmd=
case $lang in
    python)
        ext=py
        runcmd=python3
        ;;
    javascript)
        ext=js
        runcmd=node
        ;;
    go)
        ext=go
        runcmd=go run
        ;;
    *)
        echo "Wrong file type" >&2
        exit 1
        ;;
esac

# Move files
cp /tmp/code_$submission_id /tmp/$dirname/solution.$ext
cp /tmp/test_$submission_id /tmp/$dirname/test_inputs
cp /tmp/out_$submission_id /tmp/$dirname/test_outputs

# Add extra export statements for js
case $lang in
    javascript)
        echo "module.exports = { solution }" >> /tmp/$dirname/solution.$ext
        ;;
    *)
        ;;
esac

runnerfile=runner.$ext
errpath=/tmp/$dirname/errs
# Copy runner to tmp location
cp `pwd`/scripts/$runnerfile /tmp/$dirname/$runnerfile

# run command
cd /tmp/$dirname/ && timeout $RUN_LIMIT $runcmd $runnerfile 2> $errpath
stat=$?

if [[ "$stat" == "0" ]]; then
    cleanup
    exit 0
elif [[ "$stat" == "43" ]]; then
    # case when test case failed
    cat $errpath >&2
    cleanup
    exit 43
elif [[ "$stat" == "1" ]]; then
    cat $errpath >&2  # send captured error to stderr to be received by go program
    cleanup
    exit 1
elif [[ "$stat" == "124" ]]; then
    cleanup
    echo "TIMEOUT: Terminated" >&2  # send captured error to stderr to be received by go program
    exit 1
elif [[ "$stat" == "137" ]]; then
    cleanup
    echo "TIMEOUT: Killed" >&2  # send captured error to stderr to be received by go program
    exit 1
fi

echo "Unknown error occured" >&2
exit 1
