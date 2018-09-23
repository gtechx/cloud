package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"gtmsg"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"

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
var configpath = "../res/config/chatserver.config"
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
				count := int(data[0])
				chatServerMap[msg.Server] += count

				//broadcast to other chat server of uid online
				broadcastUserOnline(msg.Server, data[1:])
				// uidcount := int(data[1])
				// data = data[2:]
				// for i := 0; i < uidcount; i++ {
				// 	uid := Uint64(data[2:])

				// 	data = data[8:]
				// }
			case gtmsg.SMsgId_UserOffline:
				count := int(data[0])
				chatServerMap[msg.Server] += count

				//broadcast to other chat server of uid offline
				broadcastUserOffline(msg.Server, data[1:])
				// uidcount := int(data[1])
				// data = data[2:]
				// for i := 0; i < uidcount; i++ {
				// 	uid := Uint64(data[2:])

				// 	data = data[8:]
				// }
			}
		}

		for {
			item, err := loginServerMsgList.Pop()
			if err != nil {
				break
			}

			msg := item.(*Msg)
			//data := msg.Data
			fmt.Println("processing loginserver msg msgid " + String(msg.Msgid))
			switch msg.Msgid {
			case gtmsg.SMsgId_ReqChatServerList:
				msg.LoginConn.Write(genChatServerList())
			}
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

func broadcastUserOnline(chatserver *ChatServer, data []byte) {
	senddata := packageMsg(gtmsg.RetFrame, 0, gtmsg.SMsgId_UserOnline, data)
	for s, _ := range chatServerMap {
		if chatserver != s {
			s.conn.Write(senddata)
		}
	}
}

func broadcastUserOffline(chatserver *ChatServer, data []byte) {
	senddata := packageMsg(gtmsg.RetFrame, 0, gtmsg.SMsgId_UserOnline, data)
	for s, _ := range chatServerMap {
		if chatserver != s {
			s.conn.Write(senddata)
		}
	}
}

func packageMsg(msgtype uint8, id uint16, msgid uint16, data interface{}) []byte {
	ret := []byte{}
	databuff := Bytes(data)
	datalen := uint16(len(databuff))
	ret = append(ret, byte(msgtype))
	ret = append(ret, Bytes(id)...)
	ret = append(ret, Bytes(datalen)...)
	ret = append(ret, Bytes(msgid)...)

	if datalen > 0 {
		ret = append(ret, databuff...)
	}
	return ret
}

func readMsgHeader(conn net.Conn) (byte, uint16, uint16, uint16, []byte, error) {
	typebuff := make([]byte, 1)
	idbuff := make([]byte, 2)
	sizebuff := make([]byte, 2)
	msgidbuff := make([]byte, 2)
	var id uint16
	var size uint16
	var msgid uint16
	var databuff []byte

	_, err := io.ReadFull(conn, typebuff)
	if err != nil {
		goto end
	}

	//fmt.Println("data type:", typebuff[0])

	if typebuff[0] == gtmsg.TickFrame {
		goto end
	}

	_, err = io.ReadFull(conn, idbuff)
	if err != nil {
		goto end
	}
	id = Uint16(idbuff)

	//fmt.Println("id:", id)

	_, err = io.ReadFull(conn, sizebuff)
	if err != nil {
		goto end
	}
	size = Uint16(sizebuff)

	//fmt.Println("data size:", size)

	if size > 65535 {
		err = errors.New("too long data size")
		goto end
	}

	_, err = io.ReadFull(conn, msgidbuff)
	if err != nil {
		goto end
	}
	msgid = Uint16(msgidbuff)

	//fmt.Println("msgid:", msgid)

	if size == 0 {
		goto end
	}

	databuff = make([]byte, size)

	_, err = io.ReadFull(conn, databuff)
	if err != nil {
		goto end
	}
end:
	return typebuff[0], id, size, msgid, databuff, err
}
