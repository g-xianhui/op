package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net"
	"os"
	"runtime"
)

const (
	_ = iota
	DEBUG
	ERROR
)

func log(level int, format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func handleClient(conn net.Conn) {
	log(DEBUG, "new client[%s:%s]\n", conn.RemoteAddr(), conn.LocalAddr())
	// TODO auth process & put the data loading to right place
	accountName := "agan"
	var session uint32 = 0
	agent, err := createAgent(conn, accountName, session)
	if err != nil {
		log(ERROR, "createAgent[%s] failed: %s", accountName, err)
		return
	}
	go agentProcess(agent)
}

// database handler
var db *sql.DB
var agentcenter *AgentCenter

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	dbUser, dbPwd, dbName := "root", "", "test"
	db, err = sql.Open("mysql", dbUser+":"+dbPwd+"@/"+dbName)
	if err != nil {
		fmt.Println("Error open database: ", err)
		os.Exit(1)
	}
	defer db.Close()

	agentcenter = &AgentCenter{}
	agentcenter.init()

	l, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error listening: ", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("server started!")

	for {
		conn, err := l.Accept()
		if err != nil {
			log(ERROR, "Accept error: %s\n", err)
			continue
		}
		go handleClient(conn)
	}
}
