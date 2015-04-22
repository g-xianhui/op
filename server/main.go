package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var _ = time.Sleep

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
	// TODO auth process
	accountName := "agan"
	var session uint32 = 0
	agent := agentcenter.findByAccount(accountName)
	if agent != nil {
		mm := &InnerMsg{t: "refresh", data: &IMsgRefresh{conn: conn, session: session}}
		m := &Msg{from: 1, data: mm}
		agent.msg <- m
	} else {
		agent, err := createAgent(conn, accountName, session)
		if err != nil {
			log(ERROR, "createAgent[%s] failed: %s", accountName, err)
			conn.Close()
			return
		}
		agentcenter.addByAccount(accountName, agent)
		agent.run()
	}
}

// database handler
var db *sql.DB
var agentcenter *AgentCenter

func exit() {
	log(DEBUG, "exiting\n")
	agentcenter.exit()
	db.Close()
	os.Exit(0)
}

func sighanlder() {
	c := make(chan os.Signal)
	signal.Notify(c)
	for {
		sig := <-c
		switch sig {
		case syscall.SIGTERM, syscall.SIGINT:
			exit()
		}
		log(DEBUG, "sig[%v] catch\n", sig)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	dbUser, dbPwd, dbName := "root", "", "test"
	db, err = sql.Open("mysql", dbUser+":"+dbPwd+"@/"+dbName)
	if err != nil {
		fmt.Println("Error open database: ", err)
		os.Exit(1)
	}
	defer exit()

	agentcenter = &AgentCenter{}
	agentcenter.init()

	go sighanlder()

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
