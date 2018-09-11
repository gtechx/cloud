package main

import (
	"time"
)

func serverMonitorStart() {
	go startServerMonitor()
}

func startServerMonitor() {
	timer := time.NewTimer(time.Second * 30)

	select {
	case <-timer.C:
		dbMgr.CheckChatServerTTL()
	}
}
