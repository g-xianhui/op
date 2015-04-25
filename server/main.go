package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var _ = time.Sleep

func handleClient(conn net.Conn) {
	log(DEBUG, "new client[%s:%s]\n", conn.RemoteAddr(), conn.LocalAddr())
	// TODO auth process
	buf := make([]byte, 32)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	accountName := string(buf[:n])

	var session uint32 = 0
	agent := agentcenter.findByAccount(accountName)
	if agent != nil {
		sendInnerMsg(agent, "refresh", &IMsgRefresh{conn: conn, session: session})
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

type Env struct {
	addr                  string
	dbName, dbUser, dbPwd string
	saveinterval          int
}

var env *Env

func setupEnv(configFile string) {
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Printf("open config file failed: %s\n", err)
		os.Exit(1)
	}

	js, err := simplejson.NewFromReader(file)
	if err != nil {
		fmt.Printf("config file parse failed: %s\n", err)
		os.Exit(1)
	}

	env = &Env{}
	env.addr = js.Get("addr").MustString()
	env.dbName = js.Get("db").MustString()
	env.dbUser = js.Get("dbuser").MustString()
	env.dbPwd = js.Get("dbpwd").MustString()
	env.saveinterval = js.Get("saveinterval").MustInt()

	logLevel := js.Get("loglevel").MustString()
	logFile := js.Get("logfile").MustString()
	logInit(logFile, logLevel)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	timenow := time.Now()
	rand.Seed(timenow.Unix())

	var configFile string
	flag.StringVar(&configFile, "config", "config.json", "config file")
	flag.Parse()
	setupEnv(configFile)

	var err error
	// dbUser, dbPwd, dbName := "root", "", "test"
	db, err = sql.Open("mysql", env.dbUser+":"+env.dbPwd+"@/"+env.dbName)
	if err != nil {
		log(ERROR, "Error open database: %s\n", err)
		os.Exit(1)
	}
	defer exit()

	agentcenter = &AgentCenter{}
	agentcenter.init()

	go sighanlder()

	// l, err := net.Listen("tcp", "localhost:1234")
	l, err := net.Listen("tcp", env.addr)
	if err != nil {
		log(ERROR, "Error listening: %S\n", err)
		os.Exit(1)
	}
	defer l.Close()
	log(INFO, "server started!\n")

	for {
		conn, err := l.Accept()
		if err != nil {
			log(ERROR, "Accept error: %s\n", err)
			continue
		}
		go handleClient(conn)
	}
}
