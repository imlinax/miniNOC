package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/golang/glog"
)

const (
	CLIENT_ADDR = "0.0.0.0:24800"
	SERVER_ADDR = "0.0.0.0:24801"
)

var (
	stopChan = make(chan bool)
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
	fmt.Println("handleClient")
	for {
		buf, err := ioutil.ReadAll(cliConn)
		if err != nil || len(buf) == 0 {
			glog.Error(err)
			stopChan <- true
			return
		}
		glog.Info("cli read: ", len(buf))
		length, err := serConn.Write(buf)
		if err != nil {
			glog.Error(err)
			stopChan <- true
			return
		}
		glog.Info("server write: ", length)
	}
}

func handleServer(cliConn, serConn *net.TCPConn) {
	fmt.Println("handleServer")
	for {
		buf, err := ioutil.ReadAll(serConn)
		if err != nil || len(buf) == 0 {
			glog.Error(err)
			stopChan <- true
			return
		}
		glog.Info("server read: ", len(buf))
		length, err := cliConn.Write(buf)
		if err != nil {
			glog.Error(err)
			stopChan <- true
			return
		}
		glog.Info("cli write: ", length)
	}
}

func main() {
	cServer := createServer(CLIENT_ADDR)
	sServer := createServer(SERVER_ADDR)

	for {
		serConn, err := sServer.AcceptTCP()
		if err != nil {
			glog.Fatal(err)
		}

		cliConn, err := cServer.AcceptTCP()
		if err != nil {
			glog.Fatal(err)
		}
		serConn.Write([]byte("hello"))

		go handleClient(cliConn, serConn)
		go handleServer(cliConn, serConn)

		<-stopChan
		cliConn.Close()
		serConn.Close()

		time.Sleep(3 * time.Second)
		for len(stopChan) > 0 {
			<-stopChan
		}
	}
}
