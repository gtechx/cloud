package main

import (
	"fmt"
	"gtdb"
	"net"
	"sync"
)

var sessMap = map[uint64]map[string]ISession//{uid:{web:sess, ios:sess, android:sess}}

var count = 0

func CreateSess(conn net.Conn, tbl_appdata *gtdb.AppData, platform string) ISession {
	fmt.Println("platform:", platform)
	sess := &Sess{appdata: tbl_appdata, conn: conn, platform: platform}
	allsess := GetSess(tbl_appdata.ID)
	if sesslist == nil {
		sesslist = &sync.Map{} //map[string]ISession{}
		sesslist.Store(count, 0)
		sm.sessMap.Store(tbl_appdata.ID, sesslist)
	}
	n, _ := sesslist.Load(count)
	oldsess, ok := sesslist.Load(platform)
	sesslist.Store(platform, sess)
	sesslist.Store(count, n.(int)+1)
	//sesslist[platform] = sess
	if ok {
		oldsess.(ISession).KickOut()
	}
	return sess
}

func DelSess(sess ISession) {
	sesslist := sm.GetSess(sess.ID())
	platform := sess.Platform()
	tmpsess, ok := sesslist.Load(platform)
	//增加判断sess == sesslist[sess.Platform()],防止顶号的时候删除sess出问题
	if sesslist != nil && ok && sess == tmpsess.(ISession) {
		//delete(sesslist, sess.Platform())
		sesslist.Delete(platform)
		n, _ := sesslist.Load(count)

		if n == 1 {
			sm.sessMap.Delete(sess.ID())
		} else {
			sesslist.Store(count, n.(int)-1)
		}
	}
}

func GetSess(id uint64) *sync.Map {
	sesslist, ok := sm.sessMap.Load(id)
	if ok {
		return sesslist.(*sync.Map)
	}
	return nil
}

func SendMsgToId(id uint64, msg []byte) bool {
	sesslist := sm.GetSess(id)
	if sesslist != nil {
		flag := false
		sesslist.Range(func(key, value interface{}) bool {
			fmt.Println("sesslist.sess platform:", key)
			sess, ok := value.(ISession)
			if ok {
				tf := sess.Send(msg)
				flag = flag || tf
			}

			return true
		})
		// for _, sess := range sesslist {
		// 	fmt.Println("sesslist.sess platform:", sess.(*Sess).Platform())
		// 	tf := sess.(*Sess).Send(msg)
		// 	flag = flag || tf
		// }
		return flag
	}
	return false
}

func TrySaveOfflineMsg(id uint64, msg []byte) {
	sesslist := sm.GetSess(id)
	if sesslist == nil {
		gtdb.Manager().SendMsgToUserOffline(id, msg)
	}
}

func SetUserOnline(id uint64, platform string) uint16 {
	tbl_online := &gtdb.Online{Dataid: id, Serveraddr: srvconfig.ServerAddr, Platform: platform}
	err := gtdb.Manager().SetUserOnline(tbl_online)
	if err != nil {
		return ERR_DB
	}
	return ERR_NONE
}

func SetUserOffline(id uint64, platform string) uint16 {
	err := gtdb.Manager().SetUserOffline(id, platform)
	if err != nil {
		return ERR_DB
	}
	return ERR_NONE
}
