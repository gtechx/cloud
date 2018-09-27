package main

import (
	"fmt"
	"gtmsg"
	"io"

	"github.com/gtechx/base/collections"
	. "github.com/gtechx/base/common"
	"github.com/gtechx/base/gtnet"
)

var exchangeServerSendList = collections.NewSafeList() //*collections.SafeList
var exchangeServerClient *gtnet.Client

func messagePullStart() {
	exchangeServerClient = gtnet.NewClient("tcp", "127.0.0.1:30001", Parser)
	err := exchangeServerClient.Connect()
	if err != nil {
		panic(err.Error())
	}

	go startMessageSend()
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

func startMessageSend() {
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
		msgtype, id, size, msgid, databuff, err := readMsgHeader(reader)
		if err != nil {
			fmt.Println("Parser:" + err.Error())
			return err
		}
		fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", size, " msgid:", msgid)
		msg := &ServerEvent{Msgid: msgid, Data: databuff}
		serverEventQueue.Put(msg)
	}
	return nil
}
