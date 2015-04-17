package main

import (
	"github.com/g-xianhui/op/server/pb"
	"net"
)

type Agent struct {
	conn        net.Conn
	accountName string
	session     uint32
	// net msg
	outside chan *msg
	// msg from framework
	inner chan interface{}
	// all user data
	*pb.Role
}

func createAgent(conn net.Conn, accountName string, session uint32) (agent *Agent, err error) {
	log(DEBUG, "createAgent[%s]\n", accountName)
	agent = &Agent{conn: conn, accountName: accountName, session: session}
	agent.outside = make(chan *msg)
	agent.inner = make(chan interface{})
	if agent.Role == nil {
		agent.Role, err = loadAll(accountName)
	}
	return
}

func agentProcess(agent *Agent) {
	log(DEBUG, "agentProcess\n")
	replyRole(agent)
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
