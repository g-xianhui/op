package main

import (
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
	conn    net.Conn
	session uint32
	status  int
	init    bool
	secret  []byte
	// use for reconnect
	connecttime int64
	// use for bread recv loop
	netquitchan       chan struct{}
	msg               chan Msg
	broadcastChannels []*channelPair
	savetimer         uint
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
		if agent.savetimer == 0 {
			agent.savetimer = timemgr.AddLoop(time.Now(), time.Second*time.Duration(env.saveinterval), func() {
				sendInnerMsg(agent, "save", nil)
			})
		}
		agent.init = true
	}
	agent.setStatus(LOGINED)
}

func (agent *Agent) refresh(conn net.Conn, session uint32, secret []byte) {
	// break the old 'recv' goroutine
	if agent.getStatus() != DISCONNECTED {
		agent.netquitchan <- struct{}{}
	}

	// login on another device
	if secret != nil {
		agent.secret = secret
	}

	agent.conn = conn
	agent.session = session
	agent.setStatus(CONNECTED)
	agent.netquitchan = make(chan struct{}, 1)
	agent.connecttime = time.Now().Unix()
	// read net msg from the new connect
	go recv(agent.conn, agent.secret, agent.msg, agent.netquitchan, agent.connecttime)
}

func (agent *Agent) disconnect() {
	agent.setStatus(DISCONNECTED)
}

func (agent *Agent) save() {
	if agent.Role != nil {
		log(DEBUG, "save[%d]\n", agent.getRoleId())
		agent.Role.save()
	}
}

func createAgent(conn net.Conn, accountName string, session uint32, secret []byte) (agent *Agent, err error) {
	log(DEBUG, "createAgent[%s]\n", accountName)
	agent = &Agent{}
	if agent.account, err = loadAccount(accountName); err != nil {
		return
	}
	agent.rolelist = loadRolelist(agent.getAccountId())
	agent.conn = conn
	agent.session = session
	agent.secret = secret
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
	timemgr.DelLoop(agent.savetimer)
	agent.savetimer = 0
	agent.save()
	agent.Role = nil
}

func (agent *Agent) run() {
	agent.netquitchan = make(chan struct{}, 1)
	agent.connecttime = time.Now().Unix()
	go recv(agent.conn, agent.secret, agent.msg, agent.netquitchan, agent.connecttime)
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
