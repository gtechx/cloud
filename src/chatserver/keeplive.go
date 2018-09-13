package main

import (
	//. "github.com/gtechx/Chat/common"
	//"github.com/gtechx/base/gtnet"
	"time"
)

func keepLiveStart() {
	go startServerTTLKeep()
}

func startServerTTLKeep() {
	dbMgr.InitChatServerTTL(srvconfig.ServerAddr, 2)
	for {
		timer := time.NewTimer(time.Second * 1)

		select {
		case <-timer.C:
			dbMgr.UpdateChatServerTTL(srvconfig.ServerAddr, 2)
		}
	}
}
