package main

import "net"

type InnerMsg struct {
	t    string
	data interface{}
}

type IMsgRefresh struct {
	conn    net.Conn
	session uint32
}

func dispatchInnerMsg(agent *Agent, m *InnerMsg) {
	switch m.t {
	case "quit":
		agent.quit(SERVERCLOSE)
		done := m.data.(chan struct{})
		done <- struct{}{}
	case "refresh":
		d := m.data.(*IMsgRefresh)
		agent.refresh(d.conn, d.session)
	}
}
