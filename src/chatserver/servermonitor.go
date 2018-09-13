package main

import (
	"time"
)

func serverMonitorStart() {
	go startServerMonitor()
}

func startServerMonitor() {
	for {
		timer := time.NewTimer(time.Second * 1)

		select {
		case <-timer.C:
			checkAllServerAlive()
		}
	}
}

func checkAllServerAlive() {
	serverlist, err := dbMgr.GetChatServerList()
	if err == nil {
		for _, serveraddr := range serverlist {
			flag, err := dbMgr.IsChatServerAlive(serveraddr)
			if err != nil {
				continue
			}
			if !flag {
				dbMgr.UnRegisterChatServer(serveraddr)
			}
		}
	}
}
