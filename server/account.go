package main

import (
	"database/sql"
	"github.com/g-xianhui/op/server/pb"
	_ "github.com/go-sql-driver/mysql"
)

func newAccount(accountName string) (id uint32, err error) {
	_, err = db.Exec("insert into account(name) values(?)", accountName)
	if err != nil {
		return
	}
	err = db.QueryRow("select guid from account where name = ?", accountName).Scan(&id)
	return
}

func loadAll(accountName string) (role *pb.Role, err error) {
	var accountId uint32
	err = db.QueryRow("select guid from account where name = ?", accountName).Scan(&accountId)
	switch {
	case err == sql.ErrNoRows:
		log(DEBUG, "create new account: %s\n", accountName)
		accountId, err = newAccount(accountName)
		if err != nil {
			return
		}
	case err != nil:
		return
	}

	role = &pb.Role{}
	if role.Basic, err = loadBasic(accountId); err != nil {
		return
	}
	return
}
