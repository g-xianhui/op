package main

import (
	"database/sql"
	"github.com/g-xianhui/op/server/pb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
)

func newRole(roleid uint32) error {
	_, err := db.Exec("insert into role(guid) values(?)", roleid)
	return err
}

func loadBasic(roleid uint32) (basic *pb.RoleBasic, err error) {
	basic = &pb.RoleBasic{}
	basic.Id = proto.Uint32(roleid)

	var occ uint32
	var name string
	err = db.QueryRow("SELECT occupation, name FROM role WHERE guid = ?", roleid).Scan(&occ, &name)
	if err == sql.ErrNoRows {
		log(DEBUG, "create new user: %d\n", roleid)
		err = newRole(roleid)
	}

	if err != nil {
		return
	}

	basic.Occupation = proto.Uint32(occ)
	basic.Name = proto.String(name)
	return
}

func replyRole(agent *Agent) {
	rep := &pb.MRRoleBasic{}
	rep.Basic = agent.GetBasic()
	replyMsg(agent, pb.MRROLEBASIC, rep)
}
