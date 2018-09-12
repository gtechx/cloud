package main

import (
	"fmt"
	"time"

	. "github.com/gtechx/base/common"
)

func messagePullStart() {
	go startMessagePull()
	go startEventPull()
}

func startMessagePull() {
	for {
		data, err := dbMgr.PullOnlineMessage(srvconfig.ServerAddr)

		if err != nil {
			//fmt.Println(err.Error())
			time.Sleep(100 * time.Millisecond)
			continue
		}

		msgid := Uint16(data)
		msg := &ServerMsg{Msgid: msgid, Data: data[2:]}
		// uid := Uint64(data)
		// msg := &ServerMsg{Uid: uid, Data: data[8:]}
		serverMsgQueue.Put(msg)
		fmt.Println("put msg ", msgid, " data ", string(data[2:]))
		// if !SendMsgToId(id, data[8:]) {
		// 	dbMgr.SendMsgToUserOffline(id, data[8:])
		// }
	}
}

func startEventPull() {
	for {
		data, err := dbMgr.PullServerEvent(srvconfig.ServerAddr)

		if err != nil {
			//fmt.Println(err.Error())
			time.Sleep(100 * time.Millisecond)
			continue
		}

		msgid := Uint16(data)
		msg := &ServerEvent{Msgid: msgid, Data: data[2:]}
		serverEventQueue.Put(msg)
	}
}
