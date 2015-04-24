package main

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
)

type NetMsg struct {
	t       uint32
	session uint32
	data    []byte
}

func (m *NetMsg) getMsgType() int {
	return MSG_NET
}

func packMsg(m *NetMsg) []byte {
	l := len(m.data) + 8
	pack := make([]byte, l)
	binary.BigEndian.PutUint32(pack, m.t)
	binary.BigEndian.PutUint32(pack[4:], m.session)
	copy(pack[8:], m.data)
	return pack
}

func unpackMsg(pack []byte) *NetMsg {
	if len(pack) < 8 {
		log(ERROR, "unpackMsg failed: uncomplete package\n")
		return nil
	}
	m := &NetMsg{}
	m.t = binary.BigEndian.Uint32(pack[:4])
	m.session = binary.BigEndian.Uint32(pack[4:8])
	m.data = pack[8:]
	return m
}

type NetMsgCB func(*Agent, proto.Message)
type NetMsgHandler struct {
	p  proto.Message
	cb NetMsgCB
}

var handlers = make(map[uint32]NetMsgHandler)

func registerHandler(t uint32, p proto.Message, cb NetMsgCB) {
	handlers[t] = NetMsgHandler{p, cb}
}

func dispatchOutsideMsg(agent *Agent, m *NetMsg) {
	if agent.getStatus() == DEAD {
		return
	}

	if m.session != agent.session+1 {
		log(ERROR, "session not equal, cli[%d], svr[%d]\n", m.session, agent.session+1)
		return
	}
	agent.session++

	h, ok := handlers[m.t]
	if ok != true {
		log(ERROR, "NetMsg[%d] handler not found\n", m.t)
		return
	}

	if err := proto.Unmarshal(m.data, h.p); err != nil {
		log(ERROR, "NetMsg[%d] Unmarshal failed: %s\n", m.t, err)
		return
	}

	h.cb(agent, h.p)
}

func replyMsg(agent *Agent, t uint32, p proto.Message) {
	data, err := proto.Marshal(p)
	if err != nil {
		log(ERROR, "proto[%d] marshal failed: %s\n", t, err)
		return
	}
	m := &NetMsg{t, agent.session, data}
	if err := send(agent.conn, m); err != nil {
		log(ERROR, "role[%d] proto[%d] send failed: %s", agent.getRoleId(), m.t, err)
	}
}
