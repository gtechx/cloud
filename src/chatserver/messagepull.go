package main

import (
	"fmt"
	"gtdb"
	"time"

	. "github.com/gtechx/base/common"
)

func messagePullInit() {
	go startMessagePull()
}

func startMessagePull() {
	for {
		data, err := gtdb.Manager().PullOnlineMessage(srvconfig.ServerAddr)

		if err != nil {
			//fmt.Println(err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		id := Uint64(data[0:8])
		fmt.Println("transfer msg to ", id, " data ", string(data[8:]))
		if !SessMgr().SendMsgToId(id, data[8:]) {
			gtdb.Manager().SendMsgToUserOffline(id, data[8:])
		}
	}
}
