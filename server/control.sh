#!/bin/sh


if [ "$1" = "start" ]; then
    daemonize /root/programming/go/src/github.com/g-xianhui/op/server/server -config /root/programming/go/src/github.com/g-xianhui/op/server/config.json
elif [ "$1" = "stop" ]; then
    killall server
fi
