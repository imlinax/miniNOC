package main

import (
	"net"
	"time"

	"github.com/golang/glog"
)

const (
	CLIENT_ADDR = "0.0.0.0:24800"
	SERVER_ADDR = "0.0.0.0:24801"
)

func createServer(addr string) *net.TCPListener {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		glog.Fatal(err)
	}

	server, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		glog.Fatal(err)
	}

	return server

}
func handleClient(cliConn, serConn *net.TCPConn) {
	buf := make([]byte, 1024)
	for {
		cliConn.Read(buf)
		serConn.Write(buf)

	}
}

func handleServer(cliConn, serConn *net.TCPConn) {
	buf := make([]byte, 1024)
	for {
		serConn.Read(buf)
		cliConn.Write(buf)
	}
}

func main() {
	cServer := createServer(CLIENT_ADDR)
	sServer := createServer(SERVER_ADDR)

	serConn, err := sServer.AcceptTCP()
	if err != nil {
		glog.Fatal(err)
	}

	cliConn, err := cServer.AcceptTCP()
	if err != nil {
		glog.Fatal(err)
	}

	go handleClient(cliConn, serConn)
	go handleServer(cliConn, serConn)

	for {
		time.Sleep(30 * time.Second)
	}

}
