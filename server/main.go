package main

import (
	"fmt"
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
	agent := createAgent(conn)
	go agentProcess(agent)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("server started!")
	l, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error listening: ", err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log(ERROR, "Accept error: %s\n", err)
			continue
		}
		go handleClient(conn)
	}
}
