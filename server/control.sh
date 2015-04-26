#!/bin/sh

SVRPATH=/root/programming/go/src/github.com/g-xianhui/op/server

if [ "$1" = "start" ]; then
    killall server
    daemonize -e "$SVRPATH/err.txt" $SVRPATH/server -config $SVRPATH/config.json
elif [ "$1" = "stop" ]; then
    killall server
fi
