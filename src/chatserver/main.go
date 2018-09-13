package main

import (
	"encoding/json"
	"errors"
	"flag"
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
	"github.com/gtechx/base/gtnet"
)

var quit chan os.Signal

//var userOLMapAll = map[uint64]map[string]string{} //{uid1:{ios:serveraddr, pc:serveraddr}}
var userOLMapAll = map[uint64]map[string]bool{} //{uid1:{serveraddr1:true, serveraddr2:true}} //所有服登录用户

var sessMap = map[uint64]map[string]ISession{}              //{uid:{web:sess, ios:sess, android:sess}} //本地登录用户
var roomMapLocal = map[uint64]map[uint64]bool{}             //{rid:{uid1:true, uid2:true}} //本地登录用户所在房间 //need optimize
var uidMapAppZone = map[string]map[string]map[uint64]bool{} //{appname:{zonename:{uid1:true, uid2:true}}}

type ConnData struct {
	conn        net.Conn
	tbl_appdata *gtdb.AppData
	platform    string
}

var newConnList = collections.NewSafeList() //*collections.SafeList

type ServerEvent struct {
	Msgid uint16
	Data  []byte
}

var serverEventQueue = collections.NewSafeList() //*collections.SafeList

type ServerMsg struct {
	Msgid uint16
	Data  []byte
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
var configpath string = "../res/config/chatserver.config"

func main() {
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
	defer dbMgr.UnRegisterChatServer(srvconfig.ServerAddr)

	//keep live init
	keepLiveStart()

	//check server
	//checkAllServerAlive()

	//other server live monitor init
	serverMonitorStart()

	//read online user
	serverlist, err := dbMgr.GetChatServerList()

	if err != nil {
		panic(err.Error())
	}

	for _, serveraddr := range serverlist {
		users, err := dbMgr.GetAllOnlineUser(serveraddr)
		if err != nil {
			panic(err.Error())
		}
		for _, iuid := range users {
			uid := Uint64(iuid)
			addOLUser(serveraddr, uid)
		}
	}

	server := gtnet.NewServer()
	err = server.Start(srvconfig.ServerNet, srvconfig.ServerAddr, onNewConn)
	if err != nil {
		panic(err.Error())
	}
	defer server.Stop()
	defer dbMgr.ClearOnlineInfo(srvconfig.ServerAddr)

	//msg from other server monitor
	messagePullStart()

	fmt.Println(srvconfig.ServerNet + " server start on addr " + srvconfig.ServerAddr + " ok...")

	//frame loop
	go loop()

	<-quit

	//clear
	fmt.Println("clear...")
	server.Stop()
	dbMgr.UnRegisterChatServer(srvconfig.ServerAddr)
	dbMgr.ClearOnlineInfo(srvconfig.ServerAddr)
	dbMgr.UnInitialize()
	// var str string
	// fmt.Scanln(&str)
}

func addAppZoneUid(appname, zonename string, uid uint64) {
	zonemap, ok := uidMapAppZone[appname]
	if !ok {
		zonemap = map[string]map[uint64]bool{}
		uidMapAppZone[appname] = zonemap
	}

	sessmap, ok := zonemap[zonename]

	if !ok {
		sessmap = map[uint64]bool{}
		zonemap[zonename] = sessmap
	}
	sessmap[uid] = true
}

func removeAppZoneUid(appname, zonename string, uid uint64) {
	zonemap, ok := uidMapAppZone[appname]
	if ok {
		sessmap, ok := zonemap[zonename]

		if ok {
			delete(sessmap, uid)

			if len(sessmap) == 0 {
				delete(zonemap, zonename)

				if len(zonemap) == 0 {
					delete(uidMapAppZone, appname)
				}
			}
		}
	}
}

func addOLUser(serveraddr string, uid uint64) {
	olinfo, ok := userOLMapAll[uid]
	if !ok {
		olinfo = map[string]bool{}
		userOLMapAll[uid] = olinfo
	}
	olinfo[serveraddr] = true //other server
}

func removeOLUser(serveraddr string, uid uint64) {
	olinfo, ok := userOLMapAll[uid]
	if ok {
		_, ok := olinfo[serveraddr]
		if ok {
			delete(olinfo, serveraddr)
			if len(olinfo) == 0 {
				delete(userOLMapAll, uid)
			}
		}
	}
}

func broadcastServerEvent(msgbytes []byte) error {
	serverlist, err := dbMgr.GetChatServerList()
	if err != nil {
		return err
	}

	for _, serveraddr := range serverlist {
		if serveraddr == srvconfig.ServerAddr {
			continue
		}
		err = dbMgr.SendServerEvent(serveraddr, msgbytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func broadcastServerMsg(msgbytes []byte) error {
	serverlist, err := dbMgr.GetChatServerList()
	if err != nil {
		return err
	}

	for _, serveraddr := range serverlist {
		if serveraddr == srvconfig.ServerAddr {
			continue
		}
		err = dbMgr.SendMsgToServer(serveraddr, msgbytes)
		if err != nil {
			return err
		}
	}

	return nil
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

			//add user to userOLMapAll
			olinfo, ok := userOLMapAll[conndata.tbl_appdata.ID]
			if !ok {
				olinfo = map[string]bool{}
				userOLMapAll[conndata.tbl_appdata.ID] = olinfo
			}
			_, ok = olinfo[""]
			olinfo[""] = true //local server
			if !ok {
				//if didn't had this id logined on this server
				//get room user joined and add room info to roomMapLocal
				roomlist, err := dbMgr.GetRoomListByJoined(conndata.tbl_appdata.ID)

				if err != nil {
					fmt.Println(err.Error())
					sess.Stop()
					continue
				}

				for _, room := range roomlist {
					userlist, ok := roomMapLocal[room.Rid]
					if !ok {
						userlist = map[uint64]bool{}
						roomMapLocal[room.Rid] = userlist

						users, err := dbMgr.GetRoomUserIds(room.Rid)

						if err != nil {
							continue
						}

						for _, user := range users {
							_, ok := userOLMapAll[user.Dataid]
							userlist[user.Dataid] = ok
						}
					}
				}

				//send event to other server
				msg := &SMsgUserOnline{Uid: conndata.tbl_appdata.ID, ServerAddr: srvconfig.ServerAddr}
				msg.MsgId = SMsgId_UserOnline
				msgbytes := Bytes(msg)

				if broadcastServerEvent(msgbytes) != nil {
					fmt.Println(err.Error())
					sess.Stop()
					continue
				}

				err = dbMgr.AddOnlineUser(srvconfig.ServerAddr, conndata.tbl_appdata.ID)

				if err != nil {
					fmt.Println(err.Error())
					sess.Stop()
					continue
				}

				addAppZoneUid(conndata.tbl_appdata.Appname, conndata.tbl_appdata.Zonename, conndata.tbl_appdata.ID)
			}

			limitcount++
			dbMgr.IncrByChatServerClientCount(srvconfig.ServerAddr, 1)

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
			data := event.Data
			fmt.Println("processing server event msgid " + String(event.Msgid))
			switch event.Msgid {
			case SMsgId_ServerQuit:
				isQuit = true
			case SMsgId_UserOnline:
				uid := Uint64(data)
				serveraddr := String(data[8:])
				// strarr := strings.Split(platformandserver, "#")
				// //platform := strarr[0]
				// serveraddr := strarr[1]
				addOLUser(serveraddr, uid)
			case SMsgId_UserOffline:
				uid := Uint64(data)
				serveraddr := String(data[8:])
				removeOLUser(serveraddr, uid)
			case SMsgId_RoomAddUser:
				rid := Uint64(data)
				uid := Uint64(data[8:])
				roomusers, ok := roomMapLocal[rid]
				if ok {
					roomusers[uid] = true
				}
			case SMsgId_RoomRemoveUser:
				rid := Uint64(data)
				uid := Uint64(data[8:])
				roomusers, ok := roomMapLocal[rid]
				if ok {
					_, ok := roomusers[uid]
					if ok {
						delete(roomusers, uid)
					}
				}
			case SMsgId_RoomDimiss:
				rid := Uint64(data)
				_, ok := roomMapLocal[rid]
				if ok {
					delete(roomMapLocal, rid)
				}
			}

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
			data := msg.Data
			fmt.Println("processing server msg msgid " + String(msg.Msgid))
			switch msg.Msgid {
			case SMsgId_UserMessage:
				uid := Uint64(data)
				SendMsgToLocalUid(uid, data[8:])
			case SMsgId_RoomMessage:
				rid := Uint64(data)
				SendMsgToLocalRoom(rid, data[8:])
			case SMsgId_ZonePublicMessage:
				len := int(data[0])
				appname := String(data[1 : 1+len])
				data = data[1+len:]
				len = int(data[0])
				zonename := String(data[1 : 1+len])
				data = data[1+len:]
				SendZonePublicMsg(appname, zonename, data)
			case SMsgId_AppPublicMessage:
				len := int(data[0])
				appname := String(data[1 : 1+len])
				data = data[1+len:]
				SendAppPublicMsg(appname, data)
			case SMsgId_ServerPublicMessage:
				SendServerPublicMsg(data)
			}

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

					//只有当id对应的sesslist为0时才从userOLMapAll中删除，并且向其它服务器广播id从该服务器彻底离线的event
					//remove user from userOLMapAll
					//delete(userOLMapAll, sess.ID())
					// olinfo, ok := userOLMapAll[sess.ID()]
					// if ok {
					// 	delete(olinfo, "")
					// 	if len(olinfo) == 0 {
					// 		delete(userOLMapAll, sess.ID())
					// 	}
					// }
					removeOLUser("", sess.ID())

					//get room user joined and add room info to roomMapLocal
					roomlist, err := dbMgr.GetRoomListByJoined(sess.ID())

					if err != nil {
						continue
					}

					for _, room := range roomlist {
						userlist, ok := roomMapLocal[room.Rid]
						if ok {
							flag := true
							//check if has room use still on this server
							for uid, _ := range userlist {
								//_, ok := userOLMapAll[uid]
								_, ok := sessMap[uid]
								if ok {
									flag = false
									break
								}
							}

							if flag {
								delete(roomMapLocal, room.Rid)
							}
						}
					}

					//send event to other server
					msg := &SMsgUserOffline{Uid: sess.ID(), ServerAddr: srvconfig.ServerAddr}
					msg.MsgId = SMsgId_UserOffline
					msgbytes := Bytes(msg)

					if broadcastServerEvent(msgbytes) != nil {
						fmt.Println(err.Error())
						continue
					}

					dbMgr.RemoveOnlineUser(srvconfig.ServerAddr, sess.ID())

					removeAppZoneUid(sess.AppName(), sess.ZoneName(), sess.ID())
				}
			}

			dbMgr.IncrByChatServerClientCount(srvconfig.ServerAddr, -1)
		}

		endtime := time.Now().UnixNano()
		delta := endtime - starttime
		sleeptime := 20*1000000 - delta
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
