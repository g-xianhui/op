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
	case "rolelist":
		questrolelist(agent)
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
	case "logout":
		logout(agent)
	case "chat":
		assertParam(len(params) > 3)
		chatType, _ := strconv.ParseUint(params[1], 10, 32)
		targetId, _ := strconv.ParseUint(params[2], 10, 32)
		chat(agent, uint32(chatType), uint32(targetId), params[3])
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
	writePack(agent.conn, packMsg(m))
	agent.session++
}

func questrolelist(agent *Agent) {
	req := &pb.MQRolelist{}
	quest(agent, pb.MROLELIST, req)
}

func login(agent *Agent, roleid uint32) {
	req := &pb.MQLogin{}
	req.Roleid = proto.Uint32(roleid)
	quest(agent, pb.MLOGIN, req)
}

func createRole(agent *Agent, occ uint32, name string) {
	req := &pb.MQCreateRole{}
	req.Occ = proto.Uint32(occ)
	req.Name = proto.String(name)
	quest(agent, pb.MCREATEROLE, req)
}

func echo(agent *Agent, data []byte) {
	req := &pb.MQEcho{}
	req.Data = proto.String(string(data))
	quest(agent, pb.MECHO, req)
}

func logout(agent *Agent) {
	req := &pb.MQLogout{}
	quest(agent, pb.MLOGOUT, req)
}

func chat(agent *Agent, chatType uint32, targetId uint32, content string) {
	req := &pb.MQChat{}
	req.ChatType = proto.Uint32(chatType)
	req.TargetId = proto.Uint32(targetId)
	req.Content = proto.String(content)
	quest(agent, pb.MCHAT, req)
}
