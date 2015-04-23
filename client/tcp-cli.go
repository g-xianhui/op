package main

import (
	"fmt"
	"github.com/g-xianhui/op/client/pb"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
)

const (
	_ = iota
	DEBUG
	ERROR
)

func log(level int, format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

type Agent struct {
	conn net.Conn
	// net msg
	outside chan *msg
	inner   chan string
	session uint32
}

func createAgent(conn net.Conn) *Agent {
	log(DEBUG, "createAgent\n")
	agent := &Agent{conn: conn}
	agent.outside = make(chan *msg)
	agent.inner = make(chan string)
	return agent
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("err:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	agent := createAgent(conn)
	go recv(agent)
	go readCmd(agent)

	for {
		select {
		case m := <-agent.outside:
			dispatchOutsideMsg(agent, m)
		case cmd := <-agent.inner:
			parse(agent, cmd)
		}
	}
}

func echocb(agent *Agent, data proto.Message) {
	rep := data.(*pb.MREcho)
	log(DEBUG, "%s", rep.GetData())
}

func roleload(agent *Agent, data proto.Message) {
	rep := data.(*pb.MRRoleBasic)
	info := rep.GetBasic()
	log(DEBUG, "roleid[%d], name[%s]\n", info.GetId(), info.GetName())
}

func rolelist(agent *Agent, data proto.Message) {
	log(DEBUG, "rolelist\n")
	rep := data.(*pb.MRRolelist)
	list := rep.GetRolelist()
	for _, r := range list {
		log(DEBUG, "roleid[%d], name[%s], level[%d], occ[%d]\n", r.GetId(), r.GetName(), r.GetLevel(), r.GetOccupation())
	}
	if len(list) == 0 {
		log(DEBUG, "empty rolelist, please create a new role\n")
	}
}

func logincb(agent *Agent, data proto.Message) {
	log(DEBUG, "logincb\n")
	rep := data.(*pb.MRLogin)
	if rep.GetErrno() > 0 {
		log(DEBUG, "login failed: %d\n", rep.GetErrno())
	} else {
		log(DEBUG, "login successed\n")
	}
}

func createcb(agent *Agent, data proto.Message) {
	log(DEBUG, "createcb\n")
	rep := data.(*pb.MRCreateRole)
	errno := rep.GetErrno()
	if errno > 0 {
		log(DEBUG, "create role failed: %d\n", errno)
	} else {
		r := rep.GetBasic()
		log(DEBUG, "new role id[%d], occ[%d], level[%d], name[%s]\n", r.GetId(), r.GetOccupation(), r.GetLevel(), r.GetName())
	}
}

func init() {
	registerHandler(pb.MECHO, &pb.MREcho{}, echocb)
	registerHandler(pb.MROLELIST, &pb.MRRolelist{}, rolelist)
	registerHandler(pb.MLOGIN, &pb.MRLogin{}, logincb)
	registerHandler(pb.MCREATEROLE, &pb.MRCreateRole{}, createcb)
}
