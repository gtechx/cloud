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
	for {
		timer := time.NewTimer(time.Second * 30)

		select {
		case <-timer.C:
			dbMgr.SetChatServerTTL(srvconfig.ServerAddr, 60)
		}
	}
}
