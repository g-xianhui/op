package main

import (
	"github.com/golang/protobuf/proto"
	"net"
)

type InnerMsg struct {
	cmd string
	ud  interface{}
}

func (m *InnerMsg) getMsgType() int {
	return MSG_INNER
}

func sendInnerMsg(agent *Agent, cmd string, ud interface{}) {
	m := &InnerMsg{cmd, ud}
	agent.msg <- m
}

type InnerMsgCB func(*Agent, interface{})

func registerInnerMsgHandler(cmd string, cb InnerMsgCB) {
	innerMsgHandlers[cmd] = cb
}

func dispatchInnerMsg(agent *Agent, m *InnerMsg) {
	h, ok := innerMsgHandlers[m.cmd]
	if !ok {
		log(ERROR, "inner msg[%s] handler not found\n", m.cmd)
		return
	}
	h(agent, m.ud)
}

var innerMsgHandlers = map[string]InnerMsgCB{
	"quit":       hQuit,
	"disconnect": hDisconnect,
	"save":       hSave,
	"refresh":    hRefresh,
	"redirect":   hRedirect,
}

func hQuit(agent *Agent, ud interface{}) {
	agent.save()
	done := ud.(chan struct{})
	done <- struct{}{}
}

func hDisconnect(agent *Agent, ud interface{}) {
	agent.disconnect()
}

func hSave(agent *Agent, ud interface{}) {
	agent.save()
}

type RefreshData struct {
	conn    net.Conn
	session uint32
}

func hRefresh(agent *Agent, ud interface{}) {
	d := ud.(*RefreshData)
	agent.refresh(d.conn, d.session)
}

// NetMsg inside, send between agents and services
type NetMsgInside struct {
	t uint32
	p proto.Message
}

func hRedirect(agent *Agent, ud interface{}) {
	m := ud.(*NetMsgInside)
	replyMsg(agent, m.t, m.p)
}
