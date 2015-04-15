#!/bin/sh

# server
protoc --go_out=../server/pb/ *.proto

# client
protoc --go_out=../client/pb/ *.proto
