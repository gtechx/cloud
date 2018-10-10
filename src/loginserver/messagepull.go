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
)

var exchangeServerSendList = collections.NewSafeList() //*collections.SafeList
var exchangeServerClient *gtnet.Client

type ChatServerData struct {
	ServerAddr string `json:"serveraddr"`
	Count      int    `json:"count"`
}

var chatServerMap = map[string]int{}
var minUserServer string
var minUserCount int = 999999

func messagePullStart() {
	// exchangeServerClient = gtnet.NewClient("tcp", "127.0.0.1:30000", Parser)
	// err := exchangeServerClient.Connect()
	// if err != nil {
	// 	panic("exchangeServerClient connect error:" + err.Error())
	// }

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
	// senddata := gtmsg.PackageMsg(gtmsg.ReqFrame, 0, gtmsg.SMsgId_ReqChatServerList, nil)
	// exchangeServerClient.Send(senddata)

	for {
		timer := time.NewTimer(time.Second * 2)

		select {
		case <-timer.C:
			// senddata := gtmsg.PackageMsg(gtmsg.ReqFrame, 0, gtmsg.SMsgId_ReqChatServerList, nil)
			// exchangeServerClient.Send(senddata)
			minUserServer, _ = dbMgr.GetChatServer()
		}
	}
}

func Parser(reader io.Reader) error {
	var ok bool
	var msgtype byte
	var id uint16
	var size uint16
	var msgid uint16
	var databuff []byte
	var err error
	var mincount int
	var minserver string

	for {
		msgtype, id, size, msgid, databuff, err = gtmsg.ReadMsgHeader(reader)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", size, " msgid:", msgid)

		chatServerMap = map[string]int{}
		err = json.Unmarshal(databuff, &chatServerMap)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println("chatServerMap:", chatServerMap)
		mincount = 9999
		minserver = ""
		for saddr, count := range chatServerMap {
			if count < mincount {
				mincount = count
				minserver = saddr
			}
		}

		minUserCount, ok = chatServerMap[minUserServer]
		if !ok || mincount < minUserCount {
			minUserCount = mincount
			minUserServer = minserver
			//fmt.Println("minUserServer:", minUserServer)
		}
	}
	return nil
}
