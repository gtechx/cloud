package main

import (
	"fmt"
	"gtdb"
	"net"
)

func CreateSess(conn net.Conn, tbl_appdata *gtdb.AppData, platform string) ISession {
	fmt.Println("platform:", platform)
	sess := &Sess{appdata: tbl_appdata, conn: conn, platform: platform}
	sesslist := GetSess(tbl_appdata.ID)
	if sesslist == nil {
		sesslist = map[string]ISession{}
		sessMap[tbl_appdata.ID] = sesslist
	}
	oldsess, ok := sesslist[platform]
	sesslist[platform] = sess
	if ok {
		fmt.Println("KickOut old sess, platform:", oldsess.Platform())
		oldsess.(ISession).KickOut()
	}
	return sess
}

func DelSess(sess ISession) {
	sesslist := GetSess(sess.ID())
	platform := sess.Platform()
	tmpsess, ok := sesslist[platform]
	//增加判断sess == sesslist[sess.Platform()],防止顶号的时候删除sess出问题
	if sesslist != nil && ok && sess == tmpsess.(ISession) {
		delete(sesslist, platform)

		if len(sesslist) == 0 {
			delete(sessMap, sess.ID())
		}
	}
}

func GetSess(id uint64) map[string]ISession {
	sesslist, ok := sessMap[id]
	if ok {
		return sesslist
	}
	return nil
}

func SendMsgToLocalUid(id uint64, msg []byte) bool {
	sesslist := GetSess(id)
	if sesslist != nil {
		flag := false
		for _, sess := range sesslist {
			tf := sess.Send(msg)
			flag = flag || tf
		}

		return flag
	}
	return false
}

func SendMsgByPlatform(id uint64, platform string, msg []byte) {
	sess, _ := sessMap[id][platform]
	sess.Send(msg)
}

func SendMsgToLocalRoom(rid uint64, msg []byte) {
	roomusers, ok := roomMapLocal[rid]
	if ok {
		for uid, _ := range roomusers {
			sesslist, ok := sessMap[uid]

			if ok {
				for _, sess := range sesslist {
					sess.Send(msg)
				}
			}
		}
	}
}

func TrySaveOfflineMsg(id uint64, msg []byte) {
	sesslist := GetSess(id)
	if sesslist == nil {
		dbMgr.SendMsgToUserOffline(id, msg)
	}
}

func SetUserOnline(id uint64, platform string) uint16 {
	tbl_online := &gtdb.Online{Dataid: id, Serveraddr: srvconfig.ServerAddr, Platform: platform}
	err := dbMgr.SetUserOnline(tbl_online)
	if err != nil {
		return ERR_DB
	}
	return ERR_NONE
}

func SetUserOffline(id uint64, platform string) uint16 {
	err := dbMgr.SetUserOffline(id, platform)
	if err != nil {
		return ERR_DB
	}
	return ERR_NONE
}
