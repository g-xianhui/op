package main

import (
	"github.com/golang/protobuf/proto"
	"net"
)

type InnerMsg struct {
	cmd   string
	ud    interface{}
	reply chan interface{}
}

func (m *InnerMsg) getMsgType() int {
	return MSG_INNER
}

func packInnerMsg(cmd string, ud interface{}) *InnerMsg {
	return &InnerMsg{cmd, ud, nil}
}

func sendInnerMsg(agent *Agent, cmd string, ud interface{}) {
	m := &InnerMsg{cmd, ud, nil}
	agent.msg <- m
}

func call(agent *Agent, cmd string, ud interface{}) interface{} {
	reply := make(chan interface{})
	m := &InnerMsg{cmd, ud, reply}
	agent.msg <- m
	return <-reply
}

type InnerMsgCB func(*Agent, interface{}) interface{}

func registerInnerMsgHandler(cmd string, cb InnerMsgCB) {
	innerMsgHandlers[cmd] = cb
}

func dispatchInnerMsg(agent *Agent, m *InnerMsg) {
	h, ok := innerMsgHandlers[m.cmd]
	if !ok {
		log(ERROR, "inner msg[%s] handler not found\n", m.cmd)
		return
	}

	ret := h(agent, m.ud)
	if ret != nil {
		m.reply <- ret
	}
}

var innerMsgHandlers = map[string]InnerMsgCB{
	"quit":       hQuit,
	"disconnect": hDisconnect,
	"save":       hSave,
	"refresh":    hRefresh,
	"redirect":   hRedirect,
}

func hQuit(agent *Agent, ud interface{}) interface{} {
	agent.save()
	done := ud.(chan struct{})
	done <- struct{}{}
	return nil
}

func hDisconnect(agent *Agent, ud interface{}) interface{} {
	agent.disconnect()
	return nil
}

func hSave(agent *Agent, ud interface{}) interface{} {
	agent.save()
	return nil
}

type RefreshData struct {
	conn    net.Conn
	session uint32
}

func hRefresh(agent *Agent, ud interface{}) interface{} {
	d := ud.(*RefreshData)
	agent.refresh(d.conn, d.session)
	return nil
}

// NetMsg inside, send between agents and services
type NetMsgInside struct {
	t uint32
	p proto.Message
}

func hRedirect(agent *Agent, ud interface{}) interface{} {
	m := ud.(*NetMsgInside)
	replyMsg(agent, m.t, m.p)
	return nil
}
