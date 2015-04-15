package main

import "github.com/golang/protobuf/proto"

type msg struct {
	t       uint32
	session uint32
	data    []byte
}

type msgCB func(*Agent, proto.Message)
type msgHandler struct {
	p  proto.Message
	cb msgCB
}

var handlers = make(map[uint32]msgHandler)

func registerHandler(t uint32, p proto.Message, cb msgCB) {
	handlers[t] = msgHandler{p, cb}
}

func dispatchOutsideMsg(agent *Agent, m *msg) {
	log(DEBUG, "dispatchOutsideMsg\n")
	h, ok := handlers[m.t]
	if ok != true {
		log(ERROR, "msg[%d] handler not found\n", m.t)
		return
	}

	if err := proto.Unmarshal(m.data, h.p); err != nil {
		log(ERROR, "msg[%d] Unmarshal failed: %s\n", m.t, err)
		return
	}

	h.cb(agent, h.p)
}
