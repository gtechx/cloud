package main

import (
	"fmt"
	"net"

	//. "github.com/gtechx/base/common"
	"github.com/gtechx/base/gtnet"
)

func startLoginServerMonitor() {
	server := gtnet.NewServer()
	err := server.Start(srvconfig.ServerNet, srvconfig.ServerAddr, onLoginServerConn)
	if err != nil {
		panic(err.Error())
	}
	defer server.Stop()
}

func onLoginServerConn(conn net.Conn) {
	fmt.Println("new login conn:", conn.RemoteAddr().String())
	loginServerAddChan <- conn
	defer conn.Close()

	for {
		msgtype, id, size, msgid, databuff, err := readMsgHeader(conn)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", size, " msgid:", msgid)
		loginServerMsgList.Put(&Msg{Msgid: msgid, Data: databuff, LoginConn: conn})
	}

	loginServerRemoveChan <- conn
}
