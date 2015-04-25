package main

import (
	"math/rand"
	"net"
	"time"
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

const (
	_ = iota
	MSG_NET
	MSG_INNER
)

type Msg interface {
	getMsgType() int
}

type Agent struct {
	conn     net.Conn
	session  uint32
	status   int
	msg      chan Msg
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

	agent.setStatus(DEAD)
}

func (agent *Agent) save() {
	log(DEBUG, "save[%d]\n", agent.getRoleId())
	if agent.Role != nil {
		agent.Role.save()
	}
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
	agent.msg = make(chan Msg)
	agent.setStatus(CONNECTED)
	return
}

func timeSave(id uint32) *time.Ticker {
	// save every 5-10 minutes
	n := rand.Intn(5) + 5
	ticker := time.NewTicker(time.Minute * time.Duration(n))
	go func() {
		for _ = range ticker.C {
			if agent := agentcenter.find(id); agent != nil {
				sendInnerMsg(agent, "save", nil)
			}
		}
	}()
	return ticker
}

func (agent *Agent) run() {
	go recv(agent, agent.conn, agent.msg)
	for m := range agent.msg {
		switch m.getMsgType() {
		case MSG_NET:
			o := m.(*NetMsg)
			dispatchOutsideMsg(agent, o)
		case MSG_INNER:
			o := m.(*InnerMsg)
			dispatchInnerMsg(agent, o)
		}
	}
}
