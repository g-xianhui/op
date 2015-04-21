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
	session uint32
}

func createAgent(conn net.Conn) *Agent {
	log(DEBUG, "createAgent\n")
	agent := &Agent{conn: conn}
	agent.outside = make(chan *msg)
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

	buf := make([]byte, MAX_CLIENT_BUF)
	for {
		select {
		case m := <-agent.outside:
			dispatchOutsideMsg(agent, m)
		default:
			n, err := os.Stdin.Read(buf)
			if err != nil {
				fmt.Println("os.Stdin.Read err:", err.Error())
				return
			}
			echo(agent, buf[:n])
		}
	}

}

func echo(agent *Agent, data []byte) {
	req := &pb.MQEcho{}
	req.Data = proto.String(string(data))
	quest(agent, pb.MQECHO, req)
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

func quest(agent *Agent, t uint32, p proto.Message) {
	pack, err := proto.Marshal(p)
	if err != nil {
		log(ERROR, "Marshal failed: %s\n", err)
		return
	}
	m := &msg{}
	m.t = t
	m.session = agent.session + 1
	m.data = pack
	sendPack(agent.conn, packMsg(m))
	agent.session++
}

func login(agent *Agent, roleid uint32) {
	req := &pb.MQLogin{}
	req.Roleid = proto.Uint32(roleid)
	quest(agent, pb.MQLOGIN, req)
}

func createRole(agent *Agent, occ uint32) {
	req := &pb.MQCreateRole{}
	req.Occ = proto.Uint32(occ)
	req.Name = proto.String("agan")
	quest(agent, pb.MQCREATEROLE, req)
}

func rolelist(agent *Agent, data proto.Message) {
	log(DEBUG, "rolelist\n")
	rep := data.(*pb.MRRolelist)
	list := rep.GetRolelist()
	for _, r := range list {
		log(DEBUG, "roleid[%d], name[%s], level[%d], occ[%d]\n", r.GetId(), r.GetName(), r.GetLevel(), r.GetOccupation())
	}
	if len(list) == 0 {
		createRole(agent, 1)
	} else {
		// selectRole(agent, list[0])
		login(agent, list[0].GetId())
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
	registerHandler(pb.MRECHO, &pb.MREcho{}, echocb)
	registerHandler(pb.MRROLEBASIC, &pb.MRRoleBasic{}, roleload)
	registerHandler(pb.MRROLELIST, &pb.MRRolelist{}, rolelist)
	registerHandler(pb.MRLOGIN, &pb.MRLogin{}, logincb)
	registerHandler(pb.MRCREATEROLE, &pb.MRCreateRole{}, createcb)
}
