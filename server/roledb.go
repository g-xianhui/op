package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var _ = sql.ErrNoRows

type Role struct {
	id uint32
	// basci info is rolelist[index]
	index int
	// other task, item...
}

type RoleBasic struct {
	id         uint32
	occupation uint32
	level      uint32
	name       string
}

func loadBasic(id uint32) (basic *RoleBasic, err error) {
	basic = &RoleBasic{id: id}
	err = db.QueryRow("SELECT occupation, level, name FROM role WHERE guid = ?", id).Scan(&basic.occupation, &basic.level, &basic.name)
	return
}

func loadRolelist(accountId uint32) []*RoleBasic {
	var r1, r2, r3 uint32
	db.QueryRow("SELECT role1, role2, role3 FROM rolelist WHERE accountId = ?", accountId).Scan(&r1, &r2, &r3)
	rids := []uint32{r1, r2, r3}
	roles := make([]*RoleBasic, 0, 3)
	for _, r := range rids {
		if r > 0 {
			b, _ := loadBasic(r)
			roles = append(roles, b)
		}
	}
	return roles
}

func saveRolelist(accountId uint32, rids []uint32) {
	if len(rids) != 3 {
		return
	}
	db.Exec("DELETE FROM rolelist WHERE accountId = ?", accountId)
	db.Exec("INSERT INTO rolelist(accountId, role1, role2, role3) VALUES(?, ?, ?, ?)", accountId, rids[0], rids[1], rids[2])
}

func (role *Role) load() {
	// id := role.id
	// role.tasklist = loadTask(id)
}

func dbCreateRole(occ uint32, name string) (uint32, uint32) {
	var roleid uint32
	if _, err := db.Exec("insert into role(occupation, name) values(?, ?)", occ, name); err != nil {
		log(ERROR, "create new role failed: %s\n", err)
		return 0, ErrDBOperate
	}
	db.QueryRow("SELECT guid FROM role WHERE name = ?", name).Scan(&roleid)
	return roleid, 0
}
