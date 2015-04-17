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
	m := &msg{}
	m.t = pb.MQECHO
	m.session = agent.session + 1

	req := &pb.MQEcho{}
	req.Data = proto.String(string(data))
	pack, err := proto.Marshal(req)
	if err != nil {
		log(ERROR, "Marshal failed: %s\n", err)
		return
	}
	m.data = pack
	sendPack(agent.conn, packMsg(m))
	agent.session++
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

func init() {
	registerHandler(pb.MRECHO, &pb.MREcho{}, echocb)
	registerHandler(pb.MRROLEBASIC, &pb.MRRoleBasic{}, roleload)
}
