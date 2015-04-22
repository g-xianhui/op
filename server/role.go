package main

import (
	"github.com/g-xianhui/op/server/pb"
	"github.com/golang/protobuf/proto"
)

func toRoleBasic(r *RoleBasic) *pb.RoleBasic {
	b := &pb.RoleBasic{}
	b.Id = proto.Uint32(r.id)
	b.Occupation = proto.Uint32(r.occupation)
	b.Level = proto.Uint32(r.level)
	b.Name = proto.String(r.name)
	return b
}

func replyRolelist(agent *Agent) {
	rep := &pb.MRRolelist{}
	for _, r := range agent.rolelist {
		rep.Rolelist = append(rep.Rolelist, toRoleBasic(r))
	}
	replyMsg(agent, pb.MRROLELIST, rep)
}

func findRole(agent *Agent, id uint32) int {
	for i := range agent.rolelist {
		if agent.rolelist[i].id == id {
			return i
		}
	}
	return -1
}

func setRole(agent *Agent, i int) {
	if agent.Role != nil && agent.Role.index == i {
		return
	}
	agent.Role = &Role{id: agent.rolelist[i].id, index: i}
	agent.Role.load()
}

func login(agent *Agent, id uint32) uint32 {
	if agent.getStatus() != CONNECTED {
		return ErrLoginAtWrongStage
	}
	index := findRole(agent, id)
	if index == -1 {
		return ErrRoleNotFound
	}
	setRole(agent, index)
	agent.setStatus(LIVE)
	agentcenter.add(id, agent)
	return 0
}

func createRole(agent *Agent, occ uint32, name string) (*RoleBasic, uint32) {
	if len(agent.rolelist) >= 3 {
		return nil, ErrRolelistFull
	}
	if !agentcenter.bookName(name) {
		return nil, ErrNameAlreadyUsed
	}
	roleid, errno := dbCreateRole(occ, name)
	if errno != 0 {
		agentcenter.unbookName(name)
		return nil, errno
	} else {
		agentcenter.confirmName(name, roleid)
	}
	newrole := &RoleBasic{id: roleid, occupation: occ, name: name}
	agent.rolelist = append(agent.rolelist, newrole)
	// if crash before save rolelist, this roleid will be waste, not a big deal though
	idlist := make([]uint32, 3)
	for i, r := range agent.rolelist {
		idlist[i] = r.id
	}
	saveRolelist(agent.getAccountId(), idlist)
	return newrole, 0
}
