package main

import (
	"fmt"
	"gtmsg"
	"net"

	//. "github.com/gtechx/base/common"
	"github.com/gtechx/base/gtnet"
)

var serverForChatServer *gtnet.Server

func startChatServerMonitor() {
	serverForChatServer = gtnet.NewServer()
	err := serverForChatServer.Start(srvconfig.ServerNetForChat, srvconfig.ServerAddrForChat, onChatServerConn)
	if err != nil {
		panic(err.Error())
	}
}

func onChatServerConn(conn net.Conn) {
	fmt.Println("new chatserver conn:", conn.RemoteAddr().String())
	//check ip

	msgtype, id, size, msgid, databuff, err := gtmsg.ReadMsgHeader(conn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", size, " msgid:", msgid)
	// if msgid != 9 {
	// 	return
	// }

	chatserver := &ChatServer{conn: conn, serverAddr: string(databuff)}
	chatServerAddChan <- chatserver
	defer conn.Close()

	for {
		msgtype, id, size, msgid, databuff, err := gtmsg.ReadMsgHeader(conn)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", size, " msgid:", msgid)
		chatServerMsgList.Put(&Msg{Msgid: msgid, Data: databuff, Server: chatserver})
	}

	chatServerRemoveChan <- chatserver
}
