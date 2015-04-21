package main

import (
	"github.com/g-xianhui/op/client/pb"
	"github.com/golang/protobuf/proto"
	"os"
	"strconv"
	"strings"
)

func assertParam(b bool) {
	if !b {
		panic("params too less")
	}
}

func readCmd(agent *Agent) {
	buf := make([]byte, 256)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			log(ERROR, "os.Stdin.Read err: %s\n", err)
			break
		}
		agent.inner <- string(buf[:n])
	}
}

func parse(agent *Agent, cmdstr string) {
	defer func() {
		if err := recover(); err != nil {
			log(ERROR, "%s\n", err)
		}
	}()

	params := strings.Fields(cmdstr)
	if len(params) == 0 {
		return
	}
	cmd := params[0]
	switch cmd {
	case "echo":
		assertParam(len(params) > 1)
		echo(agent, []byte(cmdstr[5:]))
	case "login":
		assertParam(len(params) > 1)
		roleid, err := strconv.ParseUint(params[1], 10, 32)
		if err != nil {
			log(ERROR, "login roleid parse failed: %s\n", err)
		}
		login(agent, uint32(roleid))
	case "createrole":
		assertParam(len(params) > 2)
		occ, err := strconv.ParseUint(params[1], 10, 32)
		if err != nil {
			log(ERROR, "creat role occ parse failed: %s\n", err)
		}
		createRole(agent, uint32(occ), params[2])
	}
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

func createRole(agent *Agent, occ uint32, name string) {
	req := &pb.MQCreateRole{}
	req.Occ = proto.Uint32(occ)
	req.Name = proto.String(name)
	quest(agent, pb.MQCREATEROLE, req)
}

func echo(agent *Agent, data []byte) {
	req := &pb.MQEcho{}
	req.Data = proto.String(string(data))
	quest(agent, pb.MQECHO, req)
}
