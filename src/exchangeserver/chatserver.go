package main

import (
	"fmt"
	"net"

	//. "github.com/gtechx/base/common"
	"github.com/gtechx/base/gtnet"
)

func startChatServerMonitor() {
	server := gtnet.NewServer()
	err := server.Start(srvconfig.ServerNetForChat, srvconfig.ServerAddrForChat, onChatServerConn)
	if err != nil {
		panic(err.Error())
	}
	defer server.Stop()
}

func onChatServerConn(conn net.Conn) {
	fmt.Println("new conn:", conn.RemoteAddr().String())
	//check ip

	msgtype, id, size, msgid, databuff, err := readMsgHeader(conn)
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
		msgtype, id, size, msgid, databuff, err := readMsgHeader(conn)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", size, " msgid:", msgid)
		chatServerMsgList.Put(&Msg{Msgid: msgid, Data: databuff, Server: chatserver})
	}

	chatServerRemoveChan <- chatserver
}
