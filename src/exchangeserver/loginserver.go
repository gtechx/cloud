package main

import (
	"fmt"
	"net"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/base/gtnet"
)

func startLoginServerMonitor() {
	server := gtnet.NewServer()
	err = server.Start(srvconfig.ServerNet, srvconfig.ServerAddr, onLoginServerConn)
	if err != nil {
		panic(err.Error())
	}
	defer server.Stop()
}

func onLoginServerConn(conn net.Conn) {
	fmt.Println("new login conn:", conn.RemoteAddr().String())
	//defer conn.Close()
	loginServerAddChan <- conn

	for {
		msgtype, id, size, msgid, databuff, err := readMsgHeader(conn)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", size, " msgid:", msgid)
	}
}
