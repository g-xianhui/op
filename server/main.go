package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
)

func log(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Println()
}

func handleClient(conn net.Conn) {
	log("new client[%s:%s]", conn.RemoteAddr(), conn.LocalAddr())
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
			log("Accept error: %s", err)
			continue
		}
		go handleClient(conn)
	}
}
