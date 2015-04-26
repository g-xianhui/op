package main

import (
	"math/rand"
	"net"
	"time"
)

const (
	_ = iota
	CONNECTED
	LOGINED
	LOGOUT
	DISCONNECTED
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

type channelPair struct {
	c  *Broadcast
	id uint32
}

type Agent struct {
	conn              net.Conn
	session           uint32
	status            int
	init              bool
	msg               chan Msg
	broadcastChannels []*channelPair
	saveTicker        *time.Ticker
	account           *Account
	rolelist          []*RoleBasic
	// cur role data
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

func (agent *Agent) subscripte(bc *Broadcast) {
	id := subscripte(worldChannel, agent.msg)
	pair := &channelPair{bc, id}
	agent.broadcastChannels = append(agent.broadcastChannels, pair)
}

func (agent *Agent) login(id uint32) {
	if !agent.init {
		agent.broadcastChannels = []*channelPair{}
		agent.subscripte(worldChannel)
		agent.saveTicker = timeSave(id)
		agent.init = true
	}
	agent.setStatus(LOGINED)
}

func (agent *Agent) refresh(conn net.Conn, session uint32) {
	// break the old 'recv' goroutine
	if agent.getStatus() == LOGINED {
		agent.conn.Close()
	}
	agent.conn = conn
	agent.session = session
	agent.setStatus(CONNECTED)
	// read net msg from the new connect
	go recv(agent, conn, agent.msg)
}

func (agent *Agent) disconnect() {
	if agent.getStatus() == LOGINED {
		agent.setStatus(DISCONNECTED)
	}
}

func (agent *Agent) save() {
	if agent.Role != nil {
		log(DEBUG, "save[%d]\n", agent.getRoleId())
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

// release agent when memory reach threshole
func (agent *Agent) clear() {
	for _, c := range agent.broadcastChannels {
		unsubscripte(c.c, c.id)
	}
	agent.broadcastChannels = nil
	agent.saveTicker.Stop()
	agent.save()
	agent.Role = nil
}

func timeSave(id uint32) *time.Ticker {
	// save every saveinterval - saveinterval + 5 minutes
	n := rand.Intn(env.saveinterval) + 5
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
