package main

import (
	"fmt"
	"gtmsg"
	"net"

	//. "github.com/gtechx/base/common"
	"github.com/gtechx/base/gtnet"
)

var serverForLoginServer *gtnet.Server

func startLoginServerMonitor() {
	serverForLoginServer = gtnet.NewServer()
	err := serverForLoginServer.Start(srvconfig.ServerNetForLogin, srvconfig.ServerAddrForLogin, onLoginServerConn)
	if err != nil {
		panic(err.Error())
	}
}

func onLoginServerConn(conn net.Conn) {
	var msgtype byte
	var id uint16
	var size uint16
	var msgid uint16
	var databuff []byte
	var err error

	remoteaddr := conn.RemoteAddr().String()
	fmt.Println("new loginserver conn:", remoteaddr)
	loginServerAddChan <- conn
	defer conn.Close()

	for {
		msgtype, id, size, msgid, databuff, err = gtmsg.ReadMsgHeader(conn)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", size, " msgid:", msgid)
		loginServerMsgList.Put(&Msg{Msgid: msgid, Data: databuff, LoginConn: conn})
	}

	loginServerRemoveChan <- conn
	fmt.Println("loginserver:" + remoteaddr + " closed")
}
