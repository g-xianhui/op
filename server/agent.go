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
	inner    chan interface{}
	account  *Account
	rolelist []*RoleBasic
	// cur role data, use pointer for later release
	*Role
}

func (agent *Agent) getAccountId() uint32 {
	return agent.account.id
}

func (agent *Agent) getRoleId() uint32 {
	return agent.Role.id
}

func createAgent(conn net.Conn, accountName string, session uint32) (agent *Agent, err error) {
	log(DEBUG, "createAgent[%s]\n", accountName)
	agent = agentcenter.findByAccount(accountName)
	if agent == nil {
		log(DEBUG, "agent[%s] not found, try load from database\n", accountName)
		agent = &Agent{}
		agent.inner = make(chan interface{})
		agent.outside = make(chan *msg)
		if agent.account, err = loadAccount(accountName); err != nil {
			return
		}
		agent.rolelist = loadRolelist(agent.getAccountId())
		agentcenter.addByAccount(accountName, agent)
	}
	agent.conn = conn
	agent.session = session
	return
}

func agentProcess(agent *Agent) {
	log(DEBUG, "agentProcess\n")
	replyRolelist(agent)
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
