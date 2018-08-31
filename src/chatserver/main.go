package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gtdb"
	"io"
	"net"
	"os"
	"os/signal"
	"time"

	. "github.com/gtechx/base/common"
	//"github.com/gtechx/base/gtnet"

	"github.com/gtechx/base/gtnet"
)

var quit chan os.Signal

var nettype string = "tcp"
var serverAddr string = "127.0.0.1:9090"
var redisNet string = "tcp"
var redisAddr string = "192.168.93.16:6379"

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

func main() {
	//var err error
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

	// pnet := flag.String("net", "", "-net=")
	// paddr := flag.String("addr", "", "-addr=")
	// //predisnet := flag.String("redisnet", redisNet, "-redisnet=")
	// predisaddr := flag.String("redisaddr", "", "-redisaddr=")

	// flag.Parse()

	// if pnet != nil && *pnet != "" {
	// 	config.ServerNet = *pnet
	// }
	// if paddr != nil && *paddr != "" {
	// 	config.ServerAddr = *paddr
	// }
	// if predisaddr != nil && *predisaddr != "" {
	// 	config.RedisAddr = *predisaddr
	// }

	cf, err := os.Open("../res/config/chatserver.config")

	if err != nil {
		panic("can not open config file chatserver.config")
	}

	fs, err := cf.Stat()

	if err != nil {
		panic("get config file chatserver.config stat error:" + err.Error())
	}

	cbuff := make([]byte, fs.Size())
	_, err = cf.Read(cbuff)

	if err != nil {
		panic("read config file chatserver.config error:" + err.Error())
	}

	srvconfig = &serverconfig{}
	err = json.Unmarshal(cbuff, srvconfig)
	if err != nil {
		panic("json.Unmarshal config file chatserver.config error:" + err.Error())
	}

	defer gtdb.Manager().UnInitialize()
	err = gtdb.Manager().Initialize(string(cbuff))
	if err != nil {
		panic("Initialize DB err:" + err.Error())
	}

	err = gtdb.Manager().ClearOnlineInfo(srvconfig.ServerAddr)

	if err != nil {
		panic("clear online info err:" + err.Error())
	}

	//register server
	err = gtdb.Manager().RegisterServer(srvconfig.ServerAddr)

	if err != nil {
		panic("register server to gtdata.Manager err:" + err.Error())
	}
	defer gtdb.Manager().UnRegisterServer(srvconfig.ServerAddr)

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

	<-quit

	//chatServerStop()
	//gtdb.Manager().UnRegisterServer(srvconfig.ServerAddr)
	//gtdata.Manager().UnInitialize()
	//EntityManager().CleanOnlineUsers()
}

func onNewConn(conn net.Conn) {
	//EntityManager().CreateNullEntity(conn)
	fmt.Println("new conn:", conn.RemoteAddr().String())
	isok := false
	defer conn.Close()
	time.AfterFunc(15*time.Second, func() {
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
	if msgid == MsgId_ReqLogin {
		//account login
		_, ret := HandlerReqLogin(databuff)

		senddata := packageMsg(RetFrame, id, MsgId_ReqLogin, ret)
		_, err = conn.Write(senddata)

		if err != nil {
			return
		}
	} else if msgid == MsgId_ReqChatLogin {
		fmt.Println(len(databuff))
		//chat login
		buff := databuff
		slen := int(buff[0])
		account := String(buff[1 : 1+slen])
		buff = buff[1+slen:]
		slen = int(buff[0])
		password := String(buff[1 : 1+slen])
		buff = buff[1+slen:]
		slen = int(buff[0])
		appname := String(buff[1 : 1+slen])
		buff = buff[1+slen:]
		slen = int(buff[0])
		zonename := String(buff[1 : 1+slen])
		buff = buff[1+slen:]
		slen = int(buff[0])
		platform := String(buff[1 : 1+slen])

		fmt.Println(account, password, appname, zonename, platform)
		errcode, ret := HandlerReqChatLogin(account, password, appname, zonename)

		senddata := packageMsg(RetFrame, id, MsgId_ReqChatLogin, ret)
		_, err = conn.Write(senddata)

		if err != nil || errcode != ERR_NONE {
			return
		}

	waitenterchat:
		msgtype, id, size, msgid, databuff, err = readMsgHeader(conn)
		//errcode := processEnterChat(conn)
		fmt.Println(msgtype, id, size, msgid)
		if err == nil && msgid == MsgId_ReqEnterChat {
			appdataid := Uint64(databuff)
			defer SessMgr().SetUserOffline(appdataid, platform)
			//errcode, ret := HandlerReqEnterChat(appdataid)
			//errcode := ERR_NONE
			tbl_appdata, err := gtdb.Manager().GetAppData(appdataid)
			if err != nil {
				errcode = ERR_DB
			}

			errcode = SessMgr().SetUserOnline(appdataid, platform)

			senddata := packageMsg(RetFrame, id, MsgId_ReqEnterChat, errcode)
			_, err = conn.Write(senddata)

			if err != nil {
				return
			}

			if errcode == ERR_NONE {
				fmt.Println("sess start:", appdataid)
				lastremoteaddr := conn.RemoteAddr().String()
				lasttime := time.Now()
				gtdb.Manager().UpdateLastLoginInfo(appdataid, lastremoteaddr, lasttime)
				sess := SessMgr().CreateSess(conn, tbl_appdata, platform)
				sess.Start()
			}
		} else if err == nil && msgid == MsgId_ReqCreateAppdata {
			nickname := String(databuff)

			_, ret := HandlerReqCreateAppdata(appname, zonename, account, nickname, conn.RemoteAddr().String())
			senddata := packageMsg(RetFrame, id, MsgId_ReqCreateAppdata, ret)
			_, err = conn.Write(senddata)

			if err != nil {
				return
			}
			goto waitenterchat
		} else if err == nil && msgtype == TickFrame {
			goto waitenterchat
		}
	}
	fmt.Println("conn end")
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

// func processEnterChat(conn net.Conn) uint16 {
// 	isok := false
// 	time.AfterFunc(15*time.Second, func() {
// 		if !isok {
// 			conn.Close()
// 		}
// 	})

// 	msgtype, id, size, msgid, databuff, err := readMsgHeader(conn)

// 	if err != nil {
// 		return ERR_UNKNOWN
// 	}

// 	isok = true

// 	if msgid == MsgId_ReqEnterChat {
// 		appdataid := Uint64(databuff)
// 		errcode, ret := HandlerReqEnterChat(appdataid)
// 		senddata := packageMsg(RetFrame, id, MsgId_ReqEnterChat, ret)
// 		_, err = conn.Write(senddata)

// 		if err != nil {
// 			conn.Close()
// 			return ERR_UNKNOWN
// 		}
// 		return errcode
// 	} else {
// 		conn.Close()
// 		return ERR_MSG_INVALID
// 	}
// }

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
