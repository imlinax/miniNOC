package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/golang/glog"
)

const (
	ROUTER_HOST = "10.95.56.247:24801"
	APP_HOST    = "127.0.0.1:24800"
)

var (
	stopChan = make(chan bool)
)

func connect(addr string) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
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
	for {
		buf, err := ioutil.ReadAll(appConn)
		if err != nil || len(buf) == 0 {
			glog.Error(err)
			stopChan <- true
			return
		}
		fmt.Println("server read: ", len(buf))
		_, err = routerConn.Write(buf)
		if err != nil {
			glog.Error(err)
			stopChan <- true
			return
		}
	}
}

func handleRouter(appConn, routerConn *net.TCPConn) {
	for {
		buf, err := ioutil.ReadAll(routerConn)
		if err != nil || len(buf) == 0 {
			glog.Error(err)
			stopChan <- true
			return
		}
		fmt.Println("router read: ", len(buf))
		length, err := appConn.Write(buf)
		if err != nil {
			glog.Error(err)
			stopChan <- true
			return
		}
		fmt.Println("app write: ", length)
	}
}
func main() {
	for {
		routerConn := connect(ROUTER_HOST)
		fmt.Println("connect to router")
		buf := make([]byte, 1024)
		length, err := routerConn.Read(buf)
		if err != nil {
			glog.Error(err)
		}
		fmt.Println("len: ", length, "data: ", string(buf))

		appConn := connect(APP_HOST)
		fmt.Println("connect to server")
		appConn.Write(buf)

		go handleApp(appConn, routerConn)
		go handleRouter(appConn, routerConn)

		<-stopChan
		appConn.Close()
		routerConn.Close()
		time.Sleep(3 * time.Second)

		for len(stopChan) > 0 {
			<-stopChan
		}

	}

}
