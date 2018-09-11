package main

import (
	"fmt"
	"time"

	. "github.com/gtechx/base/common"
)

func messagePullInit() {
	go startMessagePull()
}

func startMessagePull() {
	for {
		data, err := dbMgr.PullOnlineMessage(srvconfig.ServerAddr)

		if err != nil {
			//fmt.Println(err.Error())
			time.Sleep(200 * time.Millisecond)
			continue
		}

		id := Uint64(data[0:8])
		fmt.Println("transfer msg to ", id, " data ", string(data[8:]))
		if !SendMsgToId(id, data[8:]) {
			dbMgr.SendMsgToUserOffline(id, data[8:])
		}
	}
}
