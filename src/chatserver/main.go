package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"gtdb"
	"gtmsg"
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

var sessMap = map[uint64]map[string]ISession{}  //{uid:{web:sess, ios:sess, android:sess}} //本地登录用户
var roomMapLocal = map[uint64]map[uint64]bool{} //{rid:{uid1:true, uid2:true}} //本地登录用户所在房间 //need optimize
var roomAdminMapLocal = map[uint64]map[uint64]bool{}
var uidMapAppZone = map[string]map[string]map[uint64]bool{} //{appname:{zonename:{uid1:true, uid2:true}}}

type UserSubData struct {
	Uid  uint64
	Data []byte
}

type UserSubEvent struct {
	Uid  uint64
	Data []byte
}

type RoomSubData struct {
	Rid  uint64
	Data []byte
}

type RoomSubEvent struct {
	Rid  uint64
	Data []byte
}

type AppSubData struct {
	Appname string
	Data    []byte
}

type AppZoneSubData struct {
	Appname  string
	Zonename string
	Data     []byte
}

var userSubMsgList = collections.NewSafeList()
var userSubEventList = collections.NewSafeList()
var roomSubMsgList = collections.NewSafeList()
var roomSubEventList = collections.NewSafeList()
var roomAdminSubMsgList = collections.NewSafeList()
var appSubMsgList = collections.NewSafeList()
var appZoneSubMsgList = collections.NewSafeList()

type ConnData struct {
	conn        net.Conn
	tbl_appdata *gtdb.AppData
	platform    string
	sessChan    chan ISession
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
var configpath = "../res/config/chatserver.config"
var configjson string

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
	//serverMonitorStart()

	//read online user
	// serverlist, err := dbMgr.GetChatServerList()

	// if err != nil {
	// 	panic(err.Error())
	// }

	// for _, serveraddr := range serverlist {
	// 	users, err := dbMgr.GetAllOnlineUser(serveraddr)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	for _, iuid := range users {
	// 		uid := Uint64(iuid)
	// 		addOLUser(serveraddr, uid)
	// 	}
	// }

	server := gtnet.NewServer()
	err = server.Start(srvconfig.ServerNet, srvconfig.ServerAddr, onNewConn)
	if err != nil {
		panic(err.Error())
	}
	defer server.Stop()
	defer dbMgr.ClearOnlineInfo(srvconfig.ServerAddr)

	//sendMsgToExchangeServer(0, []byte(srvconfig.ServerAddr))

	//msg from other server monitor
	//messagePullStart()
	dbMgr.StartPubSub(onUserSubMsg, onUserSubEvent, onRoomSubMsg, onRoomSubEvent, onRoomAdminSubMsg, onAppSubMsg, onAppZoneSubMsg)

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

func onUserSubMsg(uid uint64, data []byte) {
	userSubMsgList.Put(&UserSubData{uid, data})
}

func onUserSubEvent(uid uint64, data []byte) {
	userSubEventList.Put(&UserSubEvent{uid, data})
}

func onRoomSubMsg(rid uint64, data []byte) {
	roomSubMsgList.Put(&RoomSubData{rid, data})
}

func onRoomSubEvent(rid uint64, data []byte) {
	roomSubEventList.Put(&RoomSubEvent{rid, data})
}

func onRoomAdminSubMsg(rid uint64, data []byte) {
	roomAdminSubMsgList.Put(&RoomSubData{rid, data})
}

func onAppSubMsg(appname string, data []byte) {
	appSubMsgList.Put(&AppSubData{appname, data})
}

func onAppZoneSubMsg(appname, zonename string, data []byte) {
	appZoneSubMsgList.Put(&AppZoneSubData{appname, zonename, data})
}

func addAppZoneUid(appname, zonename string, uid uint64) {
	zonemap, ok := uidMapAppZone[appname]
	if !ok {
		zonemap = map[string]map[uint64]bool{}
		uidMapAppZone[appname] = zonemap
		dbMgr.SubAppMsg(appname)
	}

	sessmap, ok := zonemap[zonename]

	if !ok {
		sessmap = map[uint64]bool{}
		zonemap[zonename] = sessmap
		dbMgr.SubAppZoneMsg(appname, zonename)
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

				dbMgr.UnSubAppZoneMsg(appname, zonename)
				if len(zonemap) == 0 {
					delete(uidMapAppZone, appname)
					dbMgr.UnSubAppMsg(appname)
				}
			}
		}
	}
}

// func addOLUser(serveraddr string, uid uint64) {
// 	olinfo, ok := userOLMapAll[uid]
// 	if !ok {
// 		olinfo = map[string]bool{}
// 		userOLMapAll[uid] = olinfo
// 	}
// 	olinfo[serveraddr] = true //other server
// }

// func removeOLUser(serveraddr string, uid uint64) {
// 	olinfo, ok := userOLMapAll[uid]
// 	if ok {
// 		_, ok := olinfo[serveraddr]
// 		if ok {
// 			delete(olinfo, serveraddr)
// 			if len(olinfo) == 0 {
// 				delete(userOLMapAll, uid)
// 			}
// 		}
// 	}
// }

func addRoomUserToMap(rid, uid uint64) error {
	var err error
	userlist, ok := roomMapLocal[rid]
	if !ok {
		userlist = map[uint64]bool{}
		roomMapLocal[rid] = userlist

		err = dbMgr.SubRoomMsg(rid)
		if err != nil {
			return err
		}
	}
	userlist[uid] = true

	ok, err = dbMgr.IsRoomAdmin(rid, uid)
	if err != nil {
		return err
	}

	if ok {
		userlist, ok = roomAdminMapLocal[rid]
		if !ok {
			userlist = map[uint64]bool{}
			roomAdminMapLocal[rid] = userlist
		}
		userlist[uid] = true

		err = dbMgr.SubRoomAdminMsg(rid)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeRoomUserFromMap(rid, uid uint64) error {
	var err error
	userlist, ok := roomMapLocal[rid]
	if ok {
		delete(userlist, uid)

		if len(userlist) == 0 {
			delete(roomMapLocal, rid)
			err = dbMgr.UnSubRoomMsg(rid)
			if err != nil {
				return err
			}
		}
	}

	userlist, ok = roomAdminMapLocal[rid]
	if ok {
		_, ok = userlist[uid]
		if ok {
			delete(userlist, uid)

			if len(userlist) == 0 {
				delete(roomAdminMapLocal, rid)
				err = dbMgr.UnSubRoomAdminMsg(rid)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// func addRoomAdminUser(rid, uid uint64) error {
// 	err := addRoomUser(rid, uid)

// 	if err != nil {
// 		return err
// 	}

// 	err = dbMgr.SubRoomAdminMsg(rid)
// 	return err
// }

func addOLRoomUser(uid uint64) error {
	roomlist, err := dbMgr.GetRoomListByJoined(uid)

	if err != nil {
		return err
	}

	rids := []uint64{}
	adminrids := []uint64{}
	for _, room := range roomlist {
		userlist, ok := roomMapLocal[room.Rid]
		if !ok {
			userlist = map[uint64]bool{}
			roomMapLocal[room.Rid] = userlist
			rids = append(rids, room.Rid)
		}
		userlist[uid] = true

		ok, err = dbMgr.IsRoomAdmin(room.Rid, uid)
		if err != nil {
			return err
		}

		if ok {
			userlist, ok = roomAdminMapLocal[room.Rid]
			if !ok {
				userlist = map[uint64]bool{}
				roomAdminMapLocal[room.Rid] = userlist
				adminrids = append(adminrids, room.Rid)
			}
			userlist[uid] = true
		}
	}

	err = dbMgr.SubRoomMsg(rids...)
	if err != nil {
		return err
	}

	err = dbMgr.SubRoomAdminMsg(rids...)
	if err != nil {
		return err
	}

	return nil
}

func removeOLRoomUser(uid uint64) error {
	roomlist, err := dbMgr.GetRoomListByJoined(uid)

	if err != nil {
		return err
	}

	rids := []uint64{}
	adminrids := []uint64{}
	for _, room := range roomlist {
		userlist, ok := roomMapLocal[room.Rid]
		if ok {
			_, ok = userlist[uid]
			if ok {
				delete(userlist, uid)

				if len(userlist) == 0 {
					delete(roomMapLocal, room.Rid)
					rids = append(rids, room.Rid)
				}
			}
		}

		userlist, ok = roomAdminMapLocal[room.Rid]
		if ok {
			_, ok = userlist[uid]
			if ok {
				delete(userlist, uid)

				if len(userlist) == 0 {
					delete(roomAdminMapLocal, room.Rid)
					adminrids = append(adminrids, room.Rid)
				}
			}
		}
	}

	err = dbMgr.UnSubRoomMsg(rids...)
	if err != nil {
		return err
	}

	err = dbMgr.UnSubRoomAdminMsg(rids...)
	if err != nil {
		return err
	}

	return nil
}

func dismissRoom(rid uint64) {
	_, ok := roomMapLocal[rid]
	if ok {
		delete(roomMapLocal, rid)
		dbMgr.UnSubRoomMsg(rid)
	}

	_, ok = roomAdminMapLocal[rid]
	if ok {
		delete(roomAdminMapLocal, rid)
	}
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
		// uidarr := []uint64{}
		// platformarr := []string{}
		for {
			item, err := newConnList.Pop()
			if err != nil {
				break
			}

			conndata := item.(*ConnData)
			uid := conndata.tbl_appdata.ID
			_, ok := sessMap[uid]
			sess := CreateSess(conndata.conn, conndata.tbl_appdata, conndata.platform)
			sess.Start()
			conndata.sessChan <- sess

			// uidarr = append(uidarr, conndata.tbl_appdata.ID)
			// platformarr = append(platformarr, conndata.platform)

			if !ok {
				//if didn't had this id logined on this server
				//get room user joined and add room info to roomMapLocal
				if err = addOLRoomUser(uid); err != nil {
					fmt.Println(err.Error())
					sess.Stop()
					continue
				}

				addAppZoneUid(conndata.tbl_appdata.Appname, conndata.tbl_appdata.Zonename, uid)

				err = dbMgr.SubUserMsg(uid)
				if err != nil {
					fmt.Println(err.Error())
					sess.Stop()
					continue
				}
			}

			limitcount++

			if limitcount >= 100 {
				break
			}
		}

		dbMgr.IncrByChatServerClientCount(srvconfig.ServerAddr, limitcount)

		//send msg to exchange server
		// lenuid := len(uidarr)
		// if lenuid > 0 {
		// 	msg := &gtmsg.SMsgUserOnline{Uids: uidarr, Platforms: platformarr, ServerAddr: srvconfig.ServerAddr}
		// 	msgdata, _ := json.Marshal(msg)
		// 	sendMsgToExchangeServer(gtmsg.SMsgId_UserOnline, msgdata)
		// }

		limitcount = 0
		//process server event
		for {
			item, err := serverEventQueue.Pop()
			if err != nil {
				break
			}

			event := item.(*ServerEvent)
			//data := event.Data
			fmt.Println("processing server event msgid " + String(event.Msgid))
			switch event.Msgid {
			case gtmsg.SMsgId_ServerQuit:
				isQuit = true
			}

			limitcount++

			if limitcount >= 100 {
				break
			}
		}

		for {
			item, err := userSubMsgList.Pop()
			if err != nil {
				break
			}

			data := item.(*UserSubData)
			SendMsgToLocalUid(data.Uid, data.Data)
		}

		for {
			item, err := userSubEventList.Pop()
			if err != nil {
				break
			}

			data := item.(*UserSubEvent)
			eventid := Uint16(data.Data)
			switch eventid {
			case gtmsg.EventId_UserJoinRoom:
				rid := Uint64(data.Data[2:])
				addRoomUserToMap(rid, data.Uid)
			case gtmsg.EventId_UserLeaveRoom:
				rid := Uint64(data.Data[2:])
				removeRoomUserFromMap(rid, data.Uid)
			case gtmsg.EventId_UserRoomAdmin:
			case gtmsg.EventId_UserRoomUnAdmin:
			}
		}

		for {
			item, err := roomSubMsgList.Pop()
			if err != nil {
				break
			}

			data := item.(*RoomSubData)
			SendMsgToLocalRoom(data.Rid, data.Data)
		}

		for {
			item, err := roomSubEventList.Pop()
			if err != nil {
				break
			}

			data := item.(*RoomSubEvent)
			eventid := Uint16(data.Data)
			switch eventid {
			case gtmsg.EventId_RoomDismiss:
				rid := Uint64(data.Data[2:])
				dismissRoom(rid)
			}
		}

		for {
			item, err := roomAdminSubMsgList.Pop()
			if err != nil {
				break
			}

			data := item.(*RoomSubData)
			SendMsgToLocalRoomAdmin(data.Rid, data.Data)
		}

		for {
			item, err := appSubMsgList.Pop()
			if err != nil {
				break
			}

			data := item.(*AppSubData)
			SendAppPublicMsg(data.Appname, data.Data)
		}

		for {
			item, err := appZoneSubMsgList.Pop()
			if err != nil {
				break
			}

			data := item.(*AppZoneSubData)
			SendZonePublicMsg(data.Appname, data.Zonename, data.Data)
		}

		//traversal all sess, can parallel the update to diff goroutine
		for _, sesslist := range sessMap {
			for _, sess := range sesslist {
				isess := sess.(ISession)
				isess.Update()
			}
		}

		//remove sess stoped
		//uidarr = []uint64{}
		limitcount = 0
		for {
			item, err := toDeleteSessList.Pop()
			if err != nil {
				break
			}

			sess := item.(ISession)
			sesslist, ok := sessMap[sess.ID()]

			//uidarr = append(uidarr, sess.ID())

			if ok {
				delete(sesslist, sess.Platform())

				if len(sesslist) == 0 {
					delete(sessMap, sess.ID())

					//只有当id对应的sesslist为0时才从userOLMapAll中删除，并且向其它服务器广播id从该服务器彻底离线的event
					//remove user from userOLMapAll
					//removeOLUser("", sess.ID())

					//get room user joined and add room info to roomMapLocal
					if err = removeOLRoomUser(sess.ID()); err != nil {
						fmt.Println(err.Error())
						continue
					}

					//dbMgr.RemoveOnlineUser(srvconfig.ServerAddr, sess.ID())

					removeAppZoneUid(sess.AppName(), sess.ZoneName(), sess.ID())

					err = dbMgr.UnSubUserMsg(sess.ID())
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
				}
			}

			limitcount++

			if limitcount >= 100 {
				break
			}
		}

		dbMgr.IncrByChatServerClientCount(srvconfig.ServerAddr, -limitcount)
		//send event to exchange server
		// lenuid = len(uidarr)
		// if lenuid > 0 {
		// 	msg := &gtmsg.SMsgUserOffline{Uids: uidarr, ServerAddr: srvconfig.ServerAddr}
		// 	msgdata, _ := json.Marshal(msg)
		// 	sendMsgToExchangeServer(gtmsg.SMsgId_UserOffline, msgdata)
		// }

		endtime := time.Now().UnixNano()
		delta := endtime - starttime
		sleeptime := 100*1000000 - delta
		//fmt.Println("starttime:", starttime, "endtime:", endtime, " sleeptime:", sleeptime)
		if sleeptime > 0 {
			time.Sleep(time.Nanosecond * time.Duration(sleeptime))
		}
	}
}

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
		conn.Close()
		return
	}
	fmt.Println("new msg msgtype:", msgtype, " id:", id, " size:", size, " msgid:", msgid)

	sesschan := make(chan ISession, 1)
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
		newConnList.Put(&ConnData{conn, tbl_appdata, req.Platform, sesschan})
	}

	sess := <-sesschan
	sess.(*Sess).startSend()
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

func readMsgHeader(reader io.Reader) (byte, uint16, uint16, uint16, []byte, error) {
	typebuff := make([]byte, 1)
	idbuff := make([]byte, 2)
	sizebuff := make([]byte, 2)
	msgidbuff := make([]byte, 2)
	var id uint16
	var size uint16
	var msgid uint16
	var databuff []byte

	_, err := io.ReadFull(reader, typebuff)
	if err != nil {
		goto end
	}

	//fmt.Println("data type:", typebuff[0])

	if typebuff[0] == TickFrame {
		goto end
	}

	_, err = io.ReadFull(reader, idbuff)
	if err != nil {
		goto end
	}
	id = Uint16(idbuff)

	//fmt.Println("id:", id)

	_, err = io.ReadFull(reader, sizebuff)
	if err != nil {
		goto end
	}
	size = Uint16(sizebuff)

	//fmt.Println("data size:", size)

	if size > 65535 {
		err = errors.New("too long data size")
		goto end
	}

	_, err = io.ReadFull(reader, msgidbuff)
	if err != nil {
		goto end
	}
	msgid = Uint16(msgidbuff)

	//fmt.Println("msgid:", msgid)

	if size == 0 {
		goto end
	}

	databuff = make([]byte, size)

	_, err = io.ReadFull(reader, databuff)
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
