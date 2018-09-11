package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gtdb"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/gtechx/base/collections"
	. "github.com/gtechx/base/common"
	//"github.com/gtechx/base/gtnet"

	"github.com/gtechx/base/gtnet"
)

var quit chan os.Signal

var userOLMap = map[uint64]string{}        //{uid1:serveraddr, uid2:serveraddr}
var roomMap = map[uint64]map[uint64]bool{} //{rid:{uid1:true, uid2:true}}

type ConnData struct {
	conn        net.Conn
	tbl_appdata *gtdb.AppData
	platform    string
}

var newConnList = collections.NewSafeList() //*collections.SafeList

type ServerEvent struct {
	cmd  uint16
	Data []byte
}

var serverEventQueue = collections.NewSafeList() //*collections.SafeList

type ServerMsg struct {
	Uid uint64
	Msg []byte
}

var serverMsgQueue = collections.NewSafeList() //*collections.SafeList

var toDeleteSessList = collections.NewSafeList() //*collections.SafeList

type serverconfig struct {
	ServerAddr string `json:"serveraddr"`
	ServerNet  string `json:"servernet"`

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

	VerifyAddr map[string]string `jsong:"verifyaddr"`

	DefaultGroupName string `json:"defaultgroupname"`
}

var srvconfig *serverconfig
var isQuit bool
var dbMgr *gtdb.DBManager

func main() {
	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

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
	defer gtdb.Manager().UnInitialize()

	//clear online user
	err = dbMgr.ClearOnlineInfo(srvconfig.ServerAddr)

	if err != nil {
		panic("clear online info err:" + err.Error())
	}

	//register server
	err = dbMgr.RegisterChatServer(srvconfig.ServerAddr)

	if err != nil {
		panic("register server to gtdata.Manager err:" + err.Error())
	}
	defer gtdb.Manager().UnRegisterChatServer(srvconfig.ServerAddr)

	//init loadbalance
	loadBanlanceInit()

	server := gtnet.NewServer()
	err = server.Start(srvconfig.ServerNet, srvconfig.ServerAddr, onNewConn)
	if err != nil {
		panic(err.Error())
	}
	defer server.Stop()

	//keep live init
	keepLiveInit()

	//other server live monitor init
	serverMonitorInit()

	//msg from other server monitor
	messagePullInit()

	fmt.Println(srvconfig.ServerNet + " server start on addr " + srvconfig.ServerAddr + " ok...")

	//frame loop
	loop()

	<-quit

	//clear online user
	err = dbMgr.ClearOnlineInfo(srvconfig.ServerAddr)
}

func loop() {
	for {
		//check quit
		if isQuit {
			break
		}

		starttime := time.Now().UnixNano()

		limitcount := 0
		//create sess for new user
		for {
			item, err := newConnList.Pop()
			if err != nil {
				break
			}

			conndata := item.(*ConnData)
			sess := CreateSess(conndata.conn, conndata.tbl_appdata, conndata.platform)
			sess.Start()

			limitcount++

			if limitcount >= 10 {
				break
			}
		}

		limitcount = 0
		//process server event
		for {
			item, err := serverEventQueue.Pop()
			if err != nil {
				break
			}

			event := item.(*ServerEvent)

			limitcount++

			if limitcount >= 10 {
				break
			}
		}

		limitcount = 0
		//process server msg
		for {
			item, err := serverMsgQueue.Pop()
			if err != nil {
				break
			}

			msg := item.(*ServerMsg)

			limitcount++

			if limitcount >= 10 {
				break
			}
		}

		//traversal all sess, can parallel the update to diff goroutine
		for _, sesslist := range sessMap {
			for _, sess := range sesslist {
				isess := sess.(ISession)
				isess.Update()
			}
		}

		//remove sess stoped
		for {
			item, err := toDeleteSessList.Pop()
			if err != nil {
				break
			}

			sess := item.(ISession)
			sesslist, ok := sessMap[sess.ID()]

			if ok {
				delete(sesslist, sess.Platform())

				if len(sesslist) == 0 {
					delete(sessMap, sess.ID())
				}
			}
		}

		endtime := time.Now().UnixNano()
		delta := endtime - starttime
		sleeptime := 10*1000000 - delta
		if sleeptime > 0 {
			time.Sleep(time.Duration(sleeptime) * time.Nanosecond)
		}
	}
}

var configjson string

func readConfig() {
	// cf, err := os.Open("../res/config/chatserver.config")

	// if err != nil {
	// 	panic("can not open config file chatserver.config")
	// }

	// fs, err := cf.Stat()

	// if err != nil {
	// 	panic("get config file chatserver.config stat error:" + err.Error())
	// }

	// cbuff := make([]byte, fs.Size())
	// _, err = cf.Read(cbuff)

	// if err != nil {
	// 	panic("read config file chatserver.config error:" + err.Error())
	// }
	data, err := ioutil.ReadFile("../res/config/chatserver.config")
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

func onNewConn(conn net.Conn) {
	fmt.Println("new conn:", conn.RemoteAddr().String())
	isok := false
	//defer conn.Close()
	time.AfterFunc(5*time.Second, func() {
		if !isok {
			conn.Close()
		}
	})

	msgtype, id, size, msgid, databuff, err := readMsgHeader(conn)
	isok = true
	fmt.Println(msgtype, id, size, msgid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if msgid == MsgId_ReqChatLogin {
		fmt.Println(len(databuff))
		//chat login
		var errcode uint16
		var appdatabytes []byte
		tbl_appdata := &gtdb.AppData{}

		req := &MsgReqChatLogin{}
		if jsonUnMarshal(databuff, req, &errcode) {
			userdata, err := gtdb.Manager().GetChatToken(req.Token)

			if err != nil {
				errcode = ERR_DB
			} else {
				jsonUnMarshal(userdata, tbl_appdata, &errcode)
			}
		}

		ret := &MsgRetChatLogin{errcode, appdatabytes}
		senddata := packageMsg(RetFrame, id, MsgId_ReqChatLogin, ret)
		_, err = conn.Write(senddata)

		if err != nil || errcode != ERR_NONE {
			conn.Close()
			return
		}

		newConnList.Put(&ConnData{conn, tbl_appdata, req.Platform})
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
		fmt.Println(err.Error())
		goto end
	}

	//fmt.Println("data type:", typebuff[0])

	if typebuff[0] == TickFrame {
		goto end
	}

	_, err = io.ReadFull(conn, idbuff)
	if err != nil {
		fmt.Println(err.Error())
		goto end
	}
	id = Uint16(idbuff)

	//fmt.Println("id:", id)

	_, err = io.ReadFull(conn, sizebuff)
	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
		goto end
	}
end:
	return typebuff[0], id, size, msgid, databuff, err
}

//first, login with account,appname and zonename
//server will return all appdataid in the zone of app
//client need to use one of the appdataid to enter chat.

//before receive chat server chat msg, client need send ready msg to server.
//账号登录的时候发送账号、密码,返回登录成功的token
//登录聊天有两种情况
//1.聊天APP应用，没有分区
//2.游戏带分区聊天应用
//登录聊天的时候需要发送账号、密码，返回appdataidlist
//进入聊天发送appdataid, 服务器根据appdataid创建session
//客户端发送可以接受消息命令，服务器设置玩家在线
