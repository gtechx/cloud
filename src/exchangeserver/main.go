package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gtdb"

	"github.com/gtechx/base/collections"
	. "github.com/gtechx/base/common"
	"github.com/satori/go.uuid"
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
	ServerAddr string `json:"serveraddr"`
	//ServerNet  string `json:"servernet"`

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

	cf, err := os.Open("../res/config/loginserver.config")

	if err != nil {
		panic("can not open config file loginserver.config")
	}

	fs, err := cf.Stat()

	if err != nil {
		panic("get config file loginserver.config stat error:" + err.Error())
	}

	cbuff := make([]byte, fs.Size())
	_, err = cf.Read(cbuff)

	if err != nil {
		panic("read config file loginserver.config error:" + err.Error())
	}

	srvconfig = &serverconfig{}
	err = json.Unmarshal(cbuff, srvconfig)
	if err != nil {
		panic("json.Unmarshal config file loginserver.config error:" + err.Error())
	}

	dbMgr = gtdb.Manager()
	err = dbMgr.Initialize(string(cbuff))
	if err != nil {
		panic("Initialize DB err:" + err.Error())
	}
	defer dbMgr.UnInitialize()

	go startHTTPServer()
	fmt.Println("server start on addr " + srvconfig.ServerAddr + " ok...")

	go loop()

	<-quit
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
			case SMsgId_UserOnline:
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
			case SMsgId_UserOffline:
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
			data := msg.Data
			fmt.Println("processing loginserver msg msgid " + String(msg.Msgid))
			switch msg.Msgid {
			case SMsgId_ReqChatServerList:
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
	senddata := packageMsg(RetFrame, 0, SMsgId_UserOnline, data)
	for s, ucount := range chatServerMap {
		if chatserver != s {
			s.conn.Write(senddata)
		}
	}
}

func broadcastUserOffline(chatserver *ChatServer, data []byte) {
	senddata := packageMsg(RetFrame, 0, SMsgId_UserOnline, data)
	for s, ucount := range chatServerMap {
		if chatserver != s {
			s.conn.Write(senddata)
		}
	}
}

func startHTTPServer() {
	//http.HandleFunc("/serverlist", getServerList)
	http.HandleFunc("/verify", verify)
	http.HandleFunc("/login", login)

	http.HandleFunc("/serverlogin", serverlogin)

	http.HandleFunc("/chatlogin", chatlogin)
	http.HandleFunc("/chatcreateuser", chatcreateuser)
	http.ListenAndServe(":9001", nil)
}

// func getServerList(rw http.ResponseWriter, req *http.Request) {
// 	serverlist, _ := gtdb.Manager().GetServerList()

// 	ret := "{\r\n\t\"serverlist\":\r\n\t[\r\n"
// 	for i := 0; i < len(serverlist); i++ {
// 		ret += "\t\t{ \"addr\":\"" + serverlist[i] + "\" }"
// 		if i != len(serverlist)-1 {
// 			ret += ",\r\n"
// 		}
// 	}
// 	ret += "\r\n\t]\r\n"
// 	ret += "}\r\n"

// 	io.WriteString(rw, ret)
// }

type LoginRetMsg struct {
	ErrorDesc string `json:"error,omitempty"`
	ErrorCode uint16 `json:"errorcode"`
	//UID       uint64 `json:"uid,string"`
	Account string `json:"account,omitempty"`
	Token   string `json:"token,omitempty"`
}

func checkLogin(account, password string) (uint16, string) {
	if account == "" {
		return 1, "account must not null"
	}

	if password == "" {
		return 1, "password must not null"
	}

	tbl_account, err := dbMgr.GetAccount(account)

	if err != nil {
		return 3, "db error:" + err.Error()
	}

	md5password := GetSaltedPassword(password, tbl_account.Salt)
	if md5password != tbl_account.Password {
		return 4, "password not right"
	}

	return 0, ""
}

func login(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		account := req.PostFormValue("account")
		password := req.PostFormValue("password")
		var err error
		var token string
		var uu uuid.UUID

		retmsg := &LoginRetMsg{}

		retmsg.ErrorCode, retmsg.ErrorDesc = checkLogin(account, password)
		if retmsg.ErrorCode != 0 {
			goto end
		}

		uu, err = uuid.NewV4()

		if err != nil {
			retmsg.ErrorDesc = "gen uuid error:" + err.Error()
			retmsg.ErrorCode = 7
			goto end
		}

		token = uu.String()

		err = dbMgr.SaveLoginToken(account, token, srvconfig.TokenTimeout)

		if err != nil {
			retmsg.ErrorDesc = "save token error:" + err.Error()
			retmsg.ErrorCode = 8
			goto end
		}
		retmsg.Token = token
		retmsg.Account = account

	end:
		data, _ := json.Marshal(&retmsg)
		io.WriteString(rw, string(data))
	}
}

func verify(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		account := req.PostFormValue("account")
		token := req.PostFormValue("token")

		dbtoken, err := dbMgr.GetLoginToken(account)

		retmsg := &LoginRetMsg{}
		if err != nil {
			retmsg.ErrorDesc = "db error:" + err.Error()
			retmsg.ErrorCode = 3
		} else if dbtoken != token {
			retmsg.ErrorDesc = "token:" + token + " not right"
			retmsg.ErrorCode = 5
		}

		data, _ := json.Marshal(&retmsg)
		io.WriteString(rw, string(data))
	}
}

func onNewConn(conn net.Conn) {
	fmt.Println("new conn:", conn.RemoteAddr().String())
	isok := false
	//defer conn.Close()
	time.AfterFunc(5*time.Second, func() {
		if !isok {
			fmt.Println("time.AfterFunc conn close")
			conn.Close()
		}
	})

	msgtype, id, size, msgid, databuff, err := readMsgHeader(conn)
	isok = true
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", size, " msgid:", msgid)
	if msgid == MsgId_ReqChatLogin {
		//chat login
		var errcode uint16
		var appdatabytes []byte
		tbl_appdata := &gtdb.AppData{}

		req := &MsgReqChatLogin{}
		if jsonUnMarshal(databuff, req, &errcode) {
			userdata, err := dbMgr.GetChatToken(req.Token)

			if err != nil {
				errcode = ERR_DB
			} else {
				jsonUnMarshal(userdata, tbl_appdata, &errcode)
				fmt.Println("uid:", tbl_appdata.ID, " logined success")
			}
		}

		ret := &MsgRetChatLogin{errcode, appdatabytes}
		senddata := packageMsg(RetFrame, id, MsgId_ReqChatLogin, ret)
		_, err = conn.Write(senddata)

		if err != nil || errcode != ERR_NONE {
			fmt.Println(err.Error())
			conn.Close()
			return
		}

		fmt.Println(tbl_appdata)
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
		goto end
	}

	//fmt.Println("data type:", typebuff[0])

	if typebuff[0] == TickFrame {
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
