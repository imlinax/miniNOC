package main

import (
	"fmt"
	"net"
	"time"

	"github.com/golang/glog"
)

const (
	ROUTER_HOST = "127.0.0.1:24801"
	APP_HOST    = "127.0.0.1:24800"
)

func connect(addr string) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", APP_HOST)
	if err != nil {
		glog.Error(err)
	}

	Conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		glog.Error(err)
	}
	return Conn
}
func handleApp(appConn, routerConn *net.TCPConn) {
	buf := make([]byte, 1024)
	for {
		appConn.Read(buf)
		routerConn.Write(buf)
	}
}

func handleRouter(appConn, routerConn *net.TCPConn) {
	buf := make([]byte, 1024)
	for {
		routerConn.Read(buf)
		appConn.Write(buf)
	}
}
func main() {
	appConn := connect(APP_HOST)
	fmt.Println("connect to server")
	routerConn := connect(ROUTER_HOST)
	fmt.Println("connect to router")

	go handleApp(appConn, routerConn)
	go handleRouter(appConn, routerConn)

	for {
		time.Sleep(30 * time.Second)
	}

}
