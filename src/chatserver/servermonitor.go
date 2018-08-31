package main

import (
	"time"

	"gtdb"
)

func serverMonitorInit() {
	go startServerMonitor()
}

func startServerMonitor() {
	timer := time.NewTimer(time.Second * 30)

	select {
	case <-timer.C:
		gtdb.Manager().CheckServerTTL()
	}
}
