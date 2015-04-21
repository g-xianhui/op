package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Account struct {
	id   uint32
	name string
}

func createAccount(accountName string) (account *Account, err error) {
	_, err = db.Exec("insert into account(name) values(?)", accountName)
	if err != nil {
		return
	}
	account = &Account{name: accountName}
	err = db.QueryRow("select guid from account where name = ?", accountName).Scan(&account.id)
	return
}

func loadAccount(accountName string) (account *Account, err error) {
	account = &Account{name: accountName}
	err = db.QueryRow("select guid from account where name = ?", accountName).Scan(&account.id)
	switch {
	case err == sql.ErrNoRows:
		log(DEBUG, "create new account: %s\n", accountName)
		account, err = createAccount(accountName)
		if err != nil {
			return
		}
	case err != nil:
		return
	}
	return
}
