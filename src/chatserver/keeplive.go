package main

import (
	//. "github.com/gtechx/Chat/common"
	//"github.com/gtechx/base/gtnet"
	"time"

	"gtdb"
)

func keepLiveInit() {
	go startServerTTLKeep()
}

func startServerTTLKeep() {
	timer := time.NewTimer(time.Second * 30)

	select {
	case <-timer.C:
		gtdb.Manager().SetServerTTL(serverAddr, 60)
	}
}
