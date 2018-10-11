package gtdb

import (
	"strings"

	"github.com/go-redis/redis"
	. "github.com/gtechx/base/common"
)

var pubsub *redis.PubSub

func (db *DBManager) StartPubSub(usermsgcb func(uint64, []byte), usereventcb func(uint64, []byte), roommsgcb func(uint64, []byte), roomeventcb func(uint64, []byte), roomadminmsgcb func(uint64, []byte), appmsgcb func(string, []byte), appzonemsgcb func(string, string, []byte)) <-chan *redis.Message {
	pubsub = db.rd.Subscribe()
	msgchan := pubsub.Channel()
	go msgLoop(msgchan, usermsgcb, usereventcb, roommsgcb, roomeventcb, roomadminmsgcb, appmsgcb, appzonemsgcb)
	return msgchan
}

func (db *DBManager) StopPubSub() error {
	return pubsub.Close()
}

func msgLoop(msgchan <-chan *redis.Message, usermsgcb func(uint64, []byte), usereventcb func(uint64, []byte), roommsgcb func(uint64, []byte), roomeventcb func(uint64, []byte), roomadminmsgcb func(uint64, []byte), appmsgcb func(string, []byte), appzonemsgcb func(string, string, []byte)) {
	for msg := range msgchan {
		keyarr := strings.Split(msg.Channel, ":")
		count := len(keyarr)
		switch keyarr[0] {
		case "user":
			usermsgcb(Uint64(keyarr[count-1]), []byte(msg.Payload))
		case "userevent":
			usereventcb(Uint64(keyarr[count-1]), []byte(msg.Payload))
		case "room":
			roommsgcb(Uint64(keyarr[count-1]), []byte(msg.Payload))
		case "roomevent":
			roomeventcb(Uint64(keyarr[count-1]), []byte(msg.Payload))
		case "roomadmin":
			roomadminmsgcb(Uint64(keyarr[count-1]), []byte(msg.Payload))
		case "app":
			appmsgcb(String(keyarr[count-1]), []byte(msg.Payload))
		case "appzone":
			appzonemsgcb(String(keyarr[count-2]), String(keyarr[count-1]), []byte(msg.Payload))
		}
	}
}

func (db *DBManager) PubUserMsg(uid uint64, data []byte) error {
	ret := db.rd.Publish("user:msg:"+String(uid), data)
	return ret.Err()
}

func (db *DBManager) PubUserEvent(uid uint64, data []byte) error {
	ret := db.rd.Publish("userevent:msg:"+String(uid), data)
	return ret.Err()
}

func (db *DBManager) SubUserMsg(uids ...uint64) error {
	struids := []string{}
	for _, uid := range uids {
		struids = append(struids, "user:msg:"+String(uid), "userevent:msg:"+String(uid))
	}
	return pubsub.Subscribe(struids...)
}

func (db *DBManager) UnSubUserMsg(uids ...uint64) error {
	struids := []string{}
	for _, uid := range uids {
		struids = append(struids, "user:msg:"+String(uid))
	}
	return pubsub.Unsubscribe(struids...)
}

func (db *DBManager) GetUserSubNum(uid uint64) (int64, error) {
	key := "user:msg:" + String(uid)
	ret := db.rd.PubSubNumSub(key)
	err := ret.Err()

	if err != nil {
		return -1, err
	}

	return ret.Val()[key], nil
}

func (db *DBManager) PubRoomMsg(rid uint64, data []byte) error {
	ret := db.rd.Publish("room:msg:"+String(rid), data)
	return ret.Err()
}

func (db *DBManager) PubRoomEvent(rid uint64, data []byte) error {
	ret := db.rd.Publish("roomevent:msg:"+String(rid), data)
	return ret.Err()
}

func (db *DBManager) SubRoomMsg(rids ...uint64) error {
	strrids := []string{}
	for _, rid := range rids {
		strrids = append(strrids, "room:msg:"+String(rid), "roomevent:msg:"+String(rid))
	}
	return pubsub.Subscribe(strrids...)
}

func (db *DBManager) UnSubRoomMsg(rids ...uint64) error {
	strrids := []string{}
	for _, rid := range rids {
		strrids = append(strrids, "room:msg:"+String(rid))
	}
	return pubsub.Unsubscribe(strrids...)
}

func (db *DBManager) GetRoomSubNum(rid uint64) (int64, error) {
	key := "room:msg:" + String(rid)
	ret := db.rd.PubSubNumSub(key)
	err := ret.Err()

	if err != nil {
		return -1, err
	}

	return ret.Val()[key], nil
}

func (db *DBManager) PubRoomAdminMsg(rid uint64, data []byte) error {
	ret := db.rd.Publish("roomadmin:msg:"+String(rid), data)
	return ret.Err()
}

func (db *DBManager) SubRoomAdminMsg(rids ...uint64) error {
	strrids := []string{}
	for _, rid := range rids {
		strrids = append(strrids, "roomadmin:msg:"+String(rid))
	}
	return pubsub.Subscribe(strrids...)
}

func (db *DBManager) UnSubRoomAdminMsg(rids ...uint64) error {
	strrids := []string{}
	for _, rid := range rids {
		strrids = append(strrids, "roomadmin:msg:"+String(rid))
	}
	return pubsub.Unsubscribe(strrids...)
}

func (db *DBManager) GetRoomAdminSubNum(rid uint64) (int64, error) {
	key := "roomadmin:msg:" + String(rid)
	ret := db.rd.PubSubNumSub(key)
	err := ret.Err()

	if err != nil {
		return -1, err
	}

	return ret.Val()[key], nil
}

func (db *DBManager) PubAppMsg(appname string, data []byte) error {
	ret := db.rd.Publish("app:msg:"+appname, data)
	return ret.Err()
}

func (db *DBManager) SubAppMsg(appnames ...string) error {
	strappnames := []string{}
	for _, appname := range appnames {
		strappnames = append(strappnames, "app:msg:"+appname)
	}
	return pubsub.Subscribe(strappnames...)
}

func (db *DBManager) UnSubAppMsg(appnames ...string) error {
	strappnames := []string{}
	for _, appname := range appnames {
		strappnames = append(strappnames, "app:msg:"+appname)
	}
	return pubsub.Unsubscribe(strappnames...)
}

func (db *DBManager) GetAppSubNum(appname string) (int64, error) {
	key := "app:msg:" + appname
	ret := db.rd.PubSubNumSub(key)
	err := ret.Err()

	if err != nil {
		return -1, err
	}

	return ret.Val()[key], nil
}

func (db *DBManager) PubAppZoneMsg(appname, zonename string, data []byte) error {
	ret := db.rd.Publish("appzone:msg:"+appname+":"+zonename, data)
	return ret.Err()
}

func (db *DBManager) SubAppZoneMsg(appname, zonename string) error {
	return pubsub.Subscribe("appzone:msg:" + appname + ":" + zonename)
}

func (db *DBManager) UnSubAppZoneMsg(appname, zonename string) error {
	return pubsub.Unsubscribe("appzone:msg:" + appname + ":" + zonename)
}

func (db *DBManager) GetAppZoneSubNum(appname, zonename string) (int64, error) {
	key := "appzone:msg:" + appname + ":" + zonename
	ret := db.rd.PubSubNumSub(key)
	err := ret.Err()

	if err != nil {
		return -1, err
	}

	return ret.Val()[key], nil
}
