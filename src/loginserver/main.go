package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"gtdb"

	. "github.com/gtechx/base/common"
	"github.com/satori/go.uuid"
)

var quit chan os.Signal

var nettype string = "tcp"
var serverAddr string = "127.0.0.1:9090"
var redisNet string = "tcp"
var redisAddr string = "192.168.93.16:6379"

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

	messagePullStart()
	go startHTTPServer()
	fmt.Println("server start on addr " + srvconfig.ServerAddr + " ok...")

	<-quit

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
