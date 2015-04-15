package main

import (
	"net"
)

type Agent struct {
	conn    net.Conn
	session uint32
	// net msg
	outside chan *msg
	// msg from framework
	inner chan interface{}
}

func createAgent(conn net.Conn) *Agent {
	log(DEBUG, "createAgent\n")
	agent := &Agent{conn: conn}
	agent.outside = make(chan *msg)
	agent.inner = make(chan interface{})
	return agent
}

func agentProcess(agent *Agent) {
	log(DEBUG, "agentProcess\n")
	go recv(agent)
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
