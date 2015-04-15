package main

import "github.com/golang/protobuf/proto"

type msgCB func(*Agent, proto.Message)
type msgHandler struct {
	p  proto.Message
	cb msgCB
}

var handlers = make(map[uint32]msgHandler)

func registerHandler(t uint32, p proto.Message, cb msgCB) {
	handlers[t] = msgHandler{p, cb}
}
