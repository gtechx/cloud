package main

import (
	"fmt"
	"gtdb"
	"net"

	"github.com/gtechx/base/collections"
	. "github.com/gtechx/base/common"
)

type ISession interface {
	ID() uint64
	Account() string
	AppName() string
	ZoneName() string
	NickName() string
	Platform() string
	Send(buff []byte) bool
	Start()
	Stop()
	KickOut()
	Update()
	IsClosed() bool
}

type MsgData struct {
	Msgtype  byte
	Id       uint16
	Msgid    uint16
	Datasize uint16
	Data     []byte
}

type Sess struct {
	appdata  *gtdb.AppData
	conn     net.Conn
	platform string

	msgList   *collections.SafeList
	sendList  *collections.SafeList
	quitChan  chan bool
	errorChan chan bool
	isClosed  bool
}

func (s *Sess) ID() uint64 {
	return s.appdata.ID
}

func (s *Sess) Account() string {
	return s.appdata.Account
}

func (s *Sess) AppName() string {
	return s.appdata.Appname
}

func (s *Sess) ZoneName() string {
	return s.appdata.Zonename
}

func (s *Sess) NickName() string {
	return s.appdata.Nickname
}

func (s *Sess) Platform() string {
	return s.platform
}

func (s *Sess) IsClosed() bool {
	return s.isClosed
}

func (s *Sess) Start() {
	s.quitChan = make(chan bool, 2) //多个goroutine有可能会quit，所以需要两个，防止阻塞某个goroutine
	s.errorChan = make(chan bool, 1)
	s.msgList = collections.NewSafeList()
	s.sendList = collections.NewSafeList()
	go s.startRecv()
	//go s.startSend()
}

func (s *Sess) Stop() {
	s.isClosed = true
	s.quitChan <- true
}

func (s *Sess) KickOut() {
	senddata := packageMsg(RetFrame, 0, MsgId_KickOut, nil)
	s.Send(senddata)
	s.Stop()
}

func (s *Sess) Send(buff []byte) bool {
	if s.isClosed {
		return false
	}
	s.sendList.Put(buff)
	// select {
	// case s.sendChan <- buff:
	// case <-time.After(time.Millisecond * 100):
	// 	return false
	// }
	return true
}

func (s *Sess) Update() {
	limitcount := 0
	//process server msg
	for {
		item, err := s.msgList.Pop()
		if err != nil {
			break
		}

		msg := item.(*MsgData)

		errcode, ret := HandleMsg(msg.Msgid, s, msg.Data)
		if errcode == ERR_MSG_INVALID {
			fmt.Println("ERR_MSG_INVALID")
			s.Stop()
			break
		}
		if ret != nil {
			senddata := packageMsg(RetFrame, msg.Id, msg.Msgid, ret)
			s.sendList.Put(senddata)
		}

		limitcount++

		if limitcount >= 10 {
			break
		}
	}
}

var tickdata = []byte{TickFrame}

func (s *Sess) startRecv() {
	for {
		msgtype, id, datasize, msgid, databuff, err := readMsgHeader(s.conn)
		if err != nil {
			fmt.Println("readMsgHeader error:" + err.Error())
			s.errorChan <- true
			break
		} else if msgtype == TickFrame {
			s.sendList.Put(tickdata)
		}
		fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", datasize, " msgid:", msgid, " from uid ", s.ID())
		msgdata := &MsgData{msgtype, id, msgid, datasize, databuff}
		s.msgList.Put(msgdata)

		// switch msgtype {
		// case TickFrame:
		// 	//s.sendChan <- tickdata
		// 	s.sendList.Put(tickdata)
		// case EchoFrame:
		// 	senddata := packageMsg(EchoFrame, id, msgid, databuff)
		// 	//s.sendChan <- senddata
		// 	s.sendList.Put(senddata)
		// default:
		// 	fmt.Println("processing msg:", msgid, " id:", id, " datalen:", len(databuff))
		// 	if msgid != MsgId_ReqQuitChat {
		// 		errcode, ret := HandleMsg(msgid, s, databuff)
		// 		if errcode == ERR_MSG_INVALID {
		// 			fmt.Println("ERR_MSG_INVALID")
		// 			s.errorChan <- true
		// 			goto end
		// 		}
		// 		if ret != nil {
		// 			senddata := packageMsg(RetFrame, id, msgid, ret)
		// 			//s.sendChan <- senddata
		// 			s.sendList.Put(senddata)
		// 		}
		// 	} else {
		// 		s.Stop()
		// 		goto end
		// 	}
		// }
	}
	fmt.Println("session uid " + String(s.ID()) + " recv end")
}

func (s *Sess) startSend() {
	for {
		select {
		case <-s.quitChan:
			fmt.Println("sess start quit...")
			// count := len(s.sendChan)
			// for i := 0; i < count; i++ {
			// 	databuff := <-s.sendChan
			// 	_, err := s.conn.Write(databuff)
			// 	if err != nil {
			// 		fmt.Println("err Send:" + err.Error())
			// 		goto end
			// 	}
			// }
			for {
				data, err := s.sendList.Pop()
				if err != nil {
					break
				}

				databuff := data.([]byte)
				_, err = s.conn.Write(databuff)
				if err != nil {
					fmt.Println("err Send:" + err.Error())
					goto end
				}
			}
			goto end
		case <-s.errorChan:
			fmt.Println("recv s.errorChan")
			goto end
		case <-s.sendList.C:
			for {
				data, err := s.sendList.Pop()
				if err != nil {
					break
				}

				databuff := data.([]byte)

				_, err = s.conn.Write(databuff)
				if err != nil {
					fmt.Println("err Send:" + err.Error())
					// if ne, ok := err.(net.Error); ok && (ne.Temporary() || ne.Timeout()) {
					// 	//srv.logf("http: Accept error: %v; retrying in %v", err, tempDelay)
					// 	//time.Sleep(tempDelay)
					// 	continue
					// }
					goto end
				}
			}
		}
	}
end:
	//fmt.Println("remove session from sessmgr..")
	fmt.Println("session uid " + String(s.ID()) + " send end")
	//SessMgr().DelSess(s)
	//count := len(s.sendChan)
	for {
		data, err := s.sendList.Pop()
		if err != nil {
			break
		}

		databuff := data.([]byte)
		TrySaveOfflineMsg(s.ID(), databuff)
	}

	s.conn.Close()
	toDeleteSessList.Put(s)

	// for i := 0; i < count; i++ {
	// 	databuff := <-s.sendChan
	// 	SessMgr().TrySaveOfflineMsg(s.ID(), databuff)
	// }
}
