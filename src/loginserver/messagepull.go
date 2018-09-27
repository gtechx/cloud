package main

import (
	"encoding/json"
	"fmt"
	"gtmsg"
	"io"
	"time"

	"github.com/gtechx/base/collections"
	. "github.com/gtechx/base/common"
	"github.com/gtechx/base/gtnet"

	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/emirpasic/gods/utils"
)

var exchangeServerSendList = collections.NewSafeList() //*collections.SafeList
var exchangeServerClient *gtnet.Client

type ChatServerData struct {
	ServerAddr string `json:"serveraddr"`
	Count      int    `json:"count"`
}

var chatserverheap *binaryheap.Heap

func messagePullStart() {
	bhcomparator := func(a, b interface{}) int {
		csa := a.(*ChatServerData)
		csb := b.(*ChatServerData)
		fmt.Println("bhcomparator:", csa.Count, "/", csb.Count)
		return utils.IntComparator(csa.Count, csb.Count)
	}
	chatserverheap = binaryheap.NewWith(bhcomparator)

	exchangeServerClient = gtnet.NewClient("tcp", "127.0.0.1:30000", Parser)
	err := exchangeServerClient.Connect()
	if err != nil {
		panic("exchangeServerClient connect error:" + err.Error())
	}

	go startMessageSend()
	//go startEventPull()
}

func sendMsgToExchangeServer(msgid uint16, msg interface{}) {
	data := Bytes(msg) //[]byte{}
	// for _, msg := range args {
	// 	buff := Bytes(msg)
	// 	data = append(data, Bytes(uint16(len(buff)))...)
	// 	data = append(data, buff...)
	// }
	senddata := gtmsg.PackageMsg(gtmsg.ReqFrame, 0, msgid, data)
	exchangeServerSendList.Put(senddata)
}

func startMessageSend() {
	for {
		timer := time.NewTimer(time.Second * 5)

		select {
		case <-timer.C:
			senddata := gtmsg.PackageMsg(gtmsg.ReqFrame, 0, gtmsg.SMsgId_ReqChatServerList, nil)
			exchangeServerClient.Send(senddata)
		}
	}
}

func Parser(reader io.Reader) error {
	fmt.Println("start read...")
	for {
		msgtype, id, size, msgid, databuff, err := gtmsg.ReadMsgHeader(reader)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", size, " msgid:", msgid)
		// msg := &ServerEvent{Msgid: msgid, Data: databuff}
		// serverEventQueue.Put(msg)
		serverlist := map[string]int{}
		err = json.Unmarshal(databuff, &serverlist)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println("serverlist:", serverlist)
		chatserverheap.Clear()
		for saddr, count := range serverlist {
			sdata := &ChatServerData{saddr, count}
			fmt.Println("sdata:", sdata)
			chatserverheap.Push(sdata)
		}
	}
	return nil
}
