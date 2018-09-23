package main

import (
	"fmt"
	"gtmsg"
	"io"

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
		return utils.IntComparator(csa.Count, csb.Count)
	}
	chatserverheap = binaryheap.NewWith(bhcomparator)

	exchangeServerClient = gtnet.NewClient("tcp", "127.0.0.1:30000", Parser)
	err := exchangeServerClient.Connect()
	if err == nil {
		panic(err.Error())
	}

	go startMessagePull()
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

func startMessagePull() {
	defer exchangeServerClient.Close()
	for {
		select {
		case <-exchangeServerSendList.C:
			for {
				item, err := exchangeServerSendList.Pop()
				if err != nil {
					break
				}
				exchangeServerClient.Send(item.([]byte))
			}
		}
	}
}

func Parser(reader io.Reader) error {
	for {
		msgtype, id, size, msgid, databuff, err := gtmsg.ReadMsgHeader(reader)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", size, " msgid:", msgid)
		msg := &ServerEvent{Msgid: msgid, Data: databuff}
		serverEventQueue.Put(msg)
	}
	return nil
}
