package main

import (
	"github.com/golang/protobuf/proto"
	"net"
)

type Agent struct {
	conn net.Conn
	// net msg
	outside chan *msg
	// msg from framework
	inner chan interface{}
}

func createAgent(conn net.Conn) *Agent {
	log("createAgent")
	agent := &Agent{conn: conn}
	agent.outside = make(chan *msg)
	agent.inner = make(chan interface{})
	return agent
}

func agentProcess(agent *Agent) {
	log("agentProcess")
	go netio(agent)
	// TODO maybe try to break this deadloop
	for {
		select {
		case m := <-agent.outside:
			dispatchOutsideMsg(agent, m)
		case m := <-agent.inner:
			dispatchInnerMsg(agent, m)
		}
	}
}

func dispatchOutsideMsg(agent *Agent, m *msg) {
	log("dispatchOutsideMsg")
	h, ok := handlers[m.t]
	if ok != true {
		log("msg[%d] handler not found", m.t)
		return
	}

	if err := proto.Unmarshal(m.data, h.p); err != nil {
		log("msg[%d] Unmarshal failed: %s", m.t, err)
		return
	}

	h.cb(agent, h.p)
}

func dispatchInnerMsg(agent *Agent, msg interface{}) {
}
