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

func (agent *Agent) GetId() uint32 {
	return agent.Role.GetBasic().GetId()
}

func createAgent(conn net.Conn, accountName string, session uint32) (agent *Agent, err error) {
	log(DEBUG, "createAgent[%s]\n", accountName)
	agent = agentcenter.findByAccount(accountName)
	if agent == nil {
		log(DEBUG, "agent[%s] not found, try load from database\n", accountName)
		agent = &Agent{accountName: accountName}
		agent.inner = make(chan interface{})
		agent.outside = make(chan *msg)
		agent.Role, err = loadAll(accountName)
		if err != nil {
			return
		}
		agentcenter.add(accountName, agent.GetId(), agent)
	}
	agent.conn = conn
	agent.session = session
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
