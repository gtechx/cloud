package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gtmsg"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"time"

	"gtdb"

	"github.com/gtechx/base/collections"
	. "github.com/gtechx/base/common"
)

var quit chan os.Signal

var loginServerAddChan = make(chan net.Conn, 1)
var loginServerRemoveChan = make(chan net.Conn, 1)
var loginServerMap = map[net.Conn]bool{}

type ChatServer struct {
	conn       net.Conn
	serverAddr string
}

var chatServerAddChan = make(chan *ChatServer, 1)
var chatServerRemoveChan = make(chan *ChatServer, 1)
var chatServerMap = map[*ChatServer]int{}

var newConnList = collections.NewSafeList()

type Msg struct {
	Msgid     uint16
	Data      []byte
	Server    *ChatServer
	LoginConn net.Conn
}

var chatServerMsgList = collections.NewSafeList()
var loginServerMsgList = collections.NewSafeList()

type serverconfig struct {
	ServerAddrForChat  string `json:"serveraddrforchat"`
	ServerNetForChat   string `json:"servernetforchat"`
	ServerAddrForLogin string `json:"serveraddrforlogin"`
	ServerNetForLogin  string `json:"servernetforlogin"`

	RedisAddr      string `json:"redisaddr"`
	RedisPassword  string `json:"redispwd"`
	RedisDefaultDB uint64 `json:"redisdefaultdb"`
	RedisMaxConn   int    `json:"redismaxconn"`

	MysqlAddr         string `json:"mysqladdr"`
	MysqlUserPassword string `json:"mysqluserpwd"`
	MysqlDefaultDB    string `json:"mysqldefaultdb"`
	MysqlTablePrefix  string `json:"mysqltableprefix"`
	MysqlLogMode      bool   `json:"mysqllogmode"`
	MysqlMaxConn      int    `json:"mysqlmaxconn"`

	TokenTimeout int `json:"tokentimeout"`
}

var srvconfig *serverconfig
var dbMgr *gtdb.DBManager
var configpath = "../res/config/exchangeserver.config"
var configjson string

func main() {
	//var err error
	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	pconfig := flag.String("config", "", "-config=")

	flag.Parse()

	if *pconfig != "" {
		configpath = *pconfig
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			var tmp string
			fmt.Print("press enter to continue...")
			fmt.Scanln(&tmp)
		}
	}()

	readConfig()

	dbMgr = gtdb.Manager()
	err := dbMgr.Initialize(configjson)
	if err != nil {
		panic("Initialize DB err:" + err.Error())
	}
	defer dbMgr.UnInitialize()

	startLoginServerMonitor()
	startChatServerMonitor()

	fmt.Println("server start on ServerAddrForChat " + srvconfig.ServerAddrForChat + " ok...")

	go loop()

	<-quit
}

func readConfig() {
	fmt.Println("reading config:" + configpath)
	data, err := ioutil.ReadFile(configpath)
	if err != nil {
		panic("read config file chatserver.config error:" + err.Error())
	}
	configjson = string(data)

	srvconfig = &serverconfig{}
	err = json.Unmarshal(data, srvconfig)
	if err != nil {
		panic("json.Unmarshal config file chatserver.config error:" + err.Error())
	}
}

func loop() {
	for {
		starttime := time.Now().UnixNano()

		select {
		case conn := <-loginServerAddChan:
			loginServerMap[conn] = true
		case conn := <-loginServerRemoveChan:
			delete(loginServerMap, conn)
		case s := <-chatServerAddChan:
			chatServerMap[s] = 0
		case s := <-chatServerRemoveChan:
			delete(chatServerMap, s)
		default:
		}

		for {
			item, err := chatServerMsgList.Pop()
			if err != nil {
				break
			}

			msg := item.(*Msg)
			data := msg.Data
			fmt.Println("processing chatserver msg msgid " + String(msg.Msgid))
			switch msg.Msgid {
			case gtmsg.SMsgId_UserOnline:
				msgdata := &gtmsg.SMsgUserOnline{}
				err = json.Unmarshal(data, msgdata)
				if err == nil {
					chatServerMap[msg.Server] += len(msgdata.Uids)
				}

				//broadcast to other chat server of uid online
				broadcastMsgToOtherChatServer(msg.Server, gtmsg.SMsgId_UserOnline, data)
				// uidcount := int(data[1])
				// data = data[2:]
				// for i := 0; i < uidcount; i++ {
				// 	uid := Uint64(data[2:])

				// 	data = data[8:]
				// }
			case gtmsg.SMsgId_UserOffline:
				msgdata := &gtmsg.SMsgUserOffline{}
				err = json.Unmarshal(data, msgdata)
				if err == nil {
					chatServerMap[msg.Server] -= len(msgdata.Uids)
				}
			case gtmsg.SMsgId_UserMessage:
				msgdata := &gtmsg.SMsgUserMessage{}
				err = json.Unmarshal(data, msgdata)
				if err == nil {
					broadcastMsgToOtherChatServer(msg.Server, gtmsg.SMsgId_UserMessage, data)
				}
			}
		}

		for {
			item, err := loginServerMsgList.Pop()
			if err != nil {
				break
			}

			msg := item.(*Msg)
			//data := msg.Data
			switch msg.Msgid {
			case gtmsg.SMsgId_ReqChatServerList:
				retdata := genChatServerList()
				//fmt.Println("processing loginserver msg msgid " + String(msg.Msgid))
				//fmt.Println("ret data:" + string(retdata))
				senddata := gtmsg.PackageMsg(gtmsg.RetFrame, 0, gtmsg.SMsgId_ReqChatServerList, retdata)
				msg.LoginConn.Write(senddata)
			}
		}

		endtime := time.Now().UnixNano()
		delta := endtime - starttime
		sleeptime := 100*1000000 - delta
		//fmt.Println("starttime:", starttime, "endtime:", endtime, " sleeptime:", sleeptime)
		if sleeptime > 0 {
			//fmt.Println("starttime:", starttime, "endtime:", endtime, " sleeptime:", sleeptime)
			time.Sleep(time.Nanosecond * time.Duration(sleeptime))
		}
	}
}

func genChatServerList() []byte {
	serverlist := map[string]int{}
	for s, ucount := range chatServerMap {
		serverlist[s.serverAddr] = ucount
	}
	data, _ := json.Marshal(serverlist)
	return data
}

func broadcastMsgToOtherChatServer(chatserver *ChatServer, msgid uint16, data []byte) {
	senddata := gtmsg.PackageMsg(gtmsg.RetFrame, 0, msgid, data)
	for s, _ := range chatServerMap {
		if chatserver != s {
			s.conn.Write(senddata)
		}
	}
}

func broadcastUserOffline(chatserver *ChatServer, data []byte) {
	senddata := gtmsg.PackageMsg(gtmsg.RetFrame, 0, gtmsg.SMsgId_UserOnline, data)
	for s, _ := range chatServerMap {
		if chatserver != s {
			s.conn.Write(senddata)
		}
	}
}
