package main

import (
	"net"
)

const (
	_ = iota
	CONNECTED
	LIVE
	DEAD
)

const (
	_ = iota
	CLIENTDISCONNECT
	CLIENTQUIT
	ANOTHERLOGIN
	SERVERCLOSE
	KICK
)

type Msg struct {
	from int
	data interface{}
}

type Agent struct {
	conn     net.Conn
	session  uint32
	status   int
	msg      chan *Msg
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

func (agent *Agent) getStatus() int {
	return agent.status
}

func (agent *Agent) setStatus(s int) {
	agent.status = s
}

func (agent *Agent) refresh(conn net.Conn, session uint32) {
	// break the old 'recv' goroutine
	if agent.getStatus() != DEAD {
		agent.conn.Close()
	}
	agent.conn = conn
	agent.session = session
	agent.setStatus(CONNECTED)
	// read net msg from the new connect
	go recv(agent, conn, agent.msg)
}

func (agent *Agent) quit(reason int) {
	log(DEBUG, "agent[%d] quit, reason[%d]\n", agent.getRoleId(), reason)
	if agent.getStatus() == DEAD {
		return
	}

	if agent.Role != nil {
		agent.Role.save()
	}
	agent.setStatus(DEAD)
}

func createAgent(conn net.Conn, accountName string, session uint32) (agent *Agent, err error) {
	log(DEBUG, "createAgent[%s]\n", accountName)
	agent = &Agent{}
	if agent.account, err = loadAccount(accountName); err != nil {
		return
	}
	agent.rolelist = loadRolelist(agent.getAccountId())
	agent.conn = conn
	agent.session = session
	agent.msg = make(chan *Msg)
	agent.setStatus(CONNECTED)
	return
}

func (agent *Agent) run() {
	go recv(agent, agent.conn, agent.msg)
	for m := range agent.msg {
		switch m.from {
		case 0:
			o := m.data.(*NetMsg)
			dispatchOutsideMsg(agent, o)
		case 1:
			o := m.data.(*InnerMsg)
			dispatchInnerMsg(agent, o)
		}
	}
}
