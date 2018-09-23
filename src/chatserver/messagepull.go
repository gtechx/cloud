package main

import (
	"fmt"
	"gtmsg"
	"io"
	"time"

	"github.com/gtechx/base/collections"
	. "github.com/gtechx/base/common"
	"github.com/gtechx/base/gtnet"
)

var exchangeServerSendList = collections.NewSafeList() //*collections.SafeList

func messagePullStart() {
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
	senddata := packageMsg(gtmsg.ReqFrame, 0, msgid, data)
	exchangeServerSendList.Put(senddata)
}

func startMessagePull() {
	client := gtnet.NewClient("tcp", "127.0.0.1:30001", Parser)
	err := client.Connect()
	if err == nil {
		panic(err.Error())
	}
	defer client.Close()
	for {
		select {
		case <-exchangeServerSendList.C:
			for {
				item, err := exchangeServerSendList.Pop()
				if err != nil {
					break
				}
				client.Send(item.([]byte))
			}
		}
	}
	// for {
	// 	data, err := dbMgr.PullOnlineMessage(srvconfig.ServerAddr)

	// 	if err != nil {
	// 		//fmt.Println(err.Error())
	// 		time.Sleep(100 * time.Millisecond)
	// 		continue
	// 	}

	// 	msgid := Uint16(data)
	// 	msg := &ServerMsg{Msgid: msgid, Data: data[2:]}
	// 	// uid := Uint64(data)
	// 	// msg := &ServerMsg{Uid: uid, Data: data[8:]}
	// 	serverMsgQueue.Put(msg)
	// 	fmt.Println("put msg ", msgid, " data ", string(data[2:]))
	// 	// if !SendMsgToLocalUid(id, data[8:]) {
	// 	// 	dbMgr.SendMsgToUserOffline(id, data[8:])
	// 	// }
	// }
}

func startEventPull() {
	for {
		data, err := dbMgr.PullServerEvent(srvconfig.ServerAddr)

		if err != nil {
			//fmt.Println(err.Error())
			time.Sleep(100 * time.Millisecond)
			continue
		}

		msgid := Uint16(data)
		msg := &ServerEvent{Msgid: msgid, Data: data[2:]}
		serverEventQueue.Put(msg)
	}
}

func Parser(reader io.Reader) error {
	for {
		msgtype, id, size, msgid, databuff, err := readMsgHeader(reader)
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
