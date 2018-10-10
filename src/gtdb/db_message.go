package gtdb

import (
	"strings"
	"time"

	"github.com/go-redis/redis"
	. "github.com/gtechx/base/common"
)

//服务器记录所有消息？
//服务器记录7天之内的消息？如果超出，则删除更早的。msg:id:timestamp 设置定时器
//未读消息数量怎样存储

// func (db *DBManager) IsPresenceRecordExists(from, to uint64) (bool, error) {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HEXISTS", "presence:record:"+String(from), to)
// 	return redis.Bool(ret, err)
// }

// func (db *DBManager) AddPresenceRecord(from, to uint64, data []byte) error {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", "presence:record:"+String(from), to, data) //记录到发送者用户记录列表，用于校验
// 	_, err := conn.Do("HSET", "presence:"+String(to), from, data)        //记录到目的地用户presence列表
// 	return err
// }

// func (db *DBManager) RemovePresenceRecord(from, to uint64) error {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HDEL", "presence:record:"+String(from), to)
// 	return err
// }

func (db *DBManager) IsPresenceExists(id, from uint64) (int64, error) {
	// conn := db.rd.Get()
	// defer conn.Close()
	// ret, err := conn.Do("EXISTS", "user:presence:"+String(id)+":"+String(from))
	// return redis.Bool(ret, err)
	ret := db.rd.Exists("user:presence:" + String(id) + ":" + String(from))
	return ret.Result()
}

func (db *DBManager) AddPresence(from, to uint64) error {
	// conn := db.rd.Get()
	// defer conn.Close()
	// _, err := conn.Do("SET", "user:presence:"+String(to)+":"+String(from), "EX", 60*60*24*7) //记录到目的地用户presence列表
	// return err
	ret := db.rd.Set("user:presence:"+String(to)+":"+String(from), "", time.Duration(60*60*24*7)*time.Second)
	return ret.Err()
}

func (db *DBManager) RemovePresence(id, from uint64) error {
	// conn := db.rd.Get()
	// defer conn.Close()
	// _, err := conn.Do("DEL", "user:presence:"+String(id)+":"+String(from))
	// return err
	ret := db.rd.Del("user:presence:" + String(id) + ":" + String(from))
	return ret.Err()
}

// func (db *DBManager) GetAllPresence(id uint64) (map[string]string, error) {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGETALL", "user:presence:"+String(id))
// 	return redis.StringMap(ret, err) //.ByteSlices(ret, err)
// }

func (db *DBManager) AddRoomPresence(rid, appdataid uint64) error {
	// conn := db.rd.Get()
	// defer conn.Close()
	// _, err := conn.Do("SADD", "room:presence:"+String(rid), appdataid)
	// //_, err := conn.Do("SET", "roompresence:"+String(rid)+":"+String(appdataid), "", "ex 259200") //记录到目的地用户presence列表
	// return err
	ret := db.rd.SAdd("room:presence:"+String(rid), String(appdataid))
	return ret.Err()
}

func (db *DBManager) RemoveRoomPresence(rid, appdataid uint64) error {
	// conn := db.rd.Get()
	// defer conn.Close()
	// _, err := conn.Do("SREM", "room:presence:"+String(rid), appdataid)
	// //_, err := conn.Do("DEL", "roompresence:"+String(rid)+":"+String(appdataid))
	// return err
	ret := db.rd.SRem("room:presence:"+String(rid), appdataid)
	return ret.Err()
}

func (db *DBManager) GetAllRoomPresence(rid uint64) ([]uint64, error) {
	// conn := db.rd.Get()
	// defer conn.Close()
	// uids := []uint64{}

	// ret, err := conn.Do("SMEMBERS", "room:presence:"+String(rid))
	// if err != nil {
	// 	return uids, err
	// }
	// retarr, err := redis.Values(ret, err)
	// if err != nil {
	// 	return uids, err
	// }
	// err = redis.ScanSlice(retarr, &uids)
	// return uids, err //.ByteSlices(ret, err)
	ret := db.rd.SMembers("room:presence:" + String(rid))
	uids := []uint64{}
	err := ret.ScanSlice(&uids)
	return uids, err
}

func (db *DBManager) IsRoomPresenceExists(rid, appdataid uint64) (bool, error) {
	// conn := db.rd.Get()
	// defer conn.Close()
	// ret, err := conn.Do("SISMEMBER", "room:presence:"+String(rid), appdataid)
	// //ret, err := conn.Do("EXISTS", "roompresence:"+String(rid)+":"+String(appdataid))
	// return redis.Bool(ret, err)
	ret := db.rd.SIsMember("room:presence:"+String(rid), appdataid)
	return ret.Result()
}

func (db *DBManager) PullOnlineMessage(serveraddr string) ([]byte, error) {
	// conn := db.rd.Get()
	// defer conn.Close()
	// ret, err := conn.Do("LPOP", "message:"+serveraddr)
	// return redis.Bytes(ret, err)
	ret := db.rd.LPop("message:" + serveraddr)
	return ret.Bytes()
}

// func (db *DBManager) GetOfflineMessage(id uint64) ([][]byte, error) {
// 	// conn := db.rd.Get()
// 	// defer conn.Close()

// 	// ret, err := conn.Do("LRANGE", "message:offline:"+String(id), 0, -1)
// 	// datalist, err := redis.ByteSlices(ret, err)
// 	// conn.Do("LTRIM", "message:offline:"+String(id), len(datalist), -1)

// 	// return datalist, err
// 	ret := db.rd.LPop("message:"+serveraddr)
// 	return ret.Bytes()
// }

func (db *DBManager) SendMsgToServer(serveraddr string, msg []byte) error {
	// conn := db.rd.Get()
	// defer conn.Close()
	// _, err := conn.Do("RPUSH", "message:"+serveraddr, msg)
	// return err

	ret := db.rd.RPush("message:"+serveraddr, msg)
	return ret.Err()
}

func (db *DBManager) SendMsgToUserOffline(to uint64, data []byte) error {
	// conn := db.rd.Get()
	// defer conn.Close()
	// _, err := conn.Do("RPUSH", "message:offline:"+String(to), data)
	// return err
	ret := db.rd.RPush("message:offline:"+String(to), data)
	return ret.Err()
}

func (db *DBManager) AddUserMsgHistory(timestamp int64, data []byte, uids ...uint64) error {
	// conn := db.rd.Get()
	// defer conn.Close()
	// _, err := conn.Do("ZADD", "user:message:history:"+String(to), timestamp, data)
	// return err
	// ret := db.rd.ZAdd("user:message:history:"+String(to), redis.Z{Score: float64(timestamp), Member: data})
	// return ret.Err()
	pipe := db.rd.TxPipeline()
	z := redis.Z{Score: float64(timestamp), Member: data}
	for _, uid := range uids {
		pipe.ZAdd("user:message:history:"+String(uid), z)
	}
	// pipe.ZAdd("user:message:history:"+String(from), z)
	// pipe.ZAdd("user:message:history:"+String(to), z)
	//pipe.Publish("user:msg:"+String(to), data)
	_, err := pipe.Exec()
	return err
}

func (db *DBManager) GetUserMsgHistory(to uint64, mintimestamp int64) ([]string, error) {
	// conn := db.rd.Get()
	// defer conn.Close()
	// ret, err := conn.Do("ZREVRANGEBYSCORE ", "user:message:history:"+String(to), "+inf", mintimestamp, "LIMIT", 0, 100)
	// return redis.ByteSlices(ret, err)
	ret := db.rd.ZRevRangeByScore("user:message:history:"+String(to), redis.ZRangeBy{Min: "+inf", Max: String(mintimestamp), Offset: 0, Count: 100})
	return ret.Result()
}

// func (db *DBManager) AddUserPresenceHistory(to uint64, timestamp int64, data []byte) error {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("ZADD", "user:presence:history:"+String(to), timestamp, data)
// 	return err
// }

// func (db *DBManager) GetUserPresenceHistory(to uint64, mintimestamp int64) ([][]byte, error) {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("ZREVRANGEBYSCORE ", "user:presence:history:"+String(to), "+inf", mintimestamp, "LIMIT", 0, 100)
// 	return redis.ByteSlices(ret, err)
// }

func (db *DBManager) AddRoomMsgHistory(rid uint64, msg []byte, timestamp int64) error {
	// conn := db.rd.Get()
	// defer conn.Close()
	// _, err := conn.Do("ZADD", "room:message:history:"+String(rid), timestamp, msg)
	// return err
	ret := db.rd.ZAdd("room:message:history:"+String(rid), redis.Z{Score: float64(timestamp), Member: msg})
	return ret.Err()
}

func (db *DBManager) GetRoomMsgHistory(rid uint64, mintimestamp int64) ([]string, error) {
	// conn := db.rd.Get()
	// defer conn.Close()
	// ret, err := conn.Do("ZREVRANGEBYSCORE ", "room:message:history:"+String(rid), "+inf", mintimestamp, "LIMIT", 0, 100)
	// return redis.ByteSlices(ret, err)
	ret := db.rd.ZRevRangeByScore("room:message:history:"+String(rid), redis.ZRangeBy{Min: "+inf", Max: String(mintimestamp), Offset: 0, Count: 100})
	return ret.Result()
}

// func (db *DBManager) AddUserToRoomApplyList(rid, uid uint64, timestamp int64) error {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("ZADD", "room:apply:list:"+String(rid), timestamp, uid)
// 	return err
// }

// func (db *DBManager) RemoveUserFromRoomApplyList(rid, uid uint64) error {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("ZREM", "room:apply:list:"+String(rid), uid)
// 	return err
// }

func (db *DBManager) SaveLoginToken(account, token string, timeout int) error {
	// conn := db.rd.Get()
	// defer conn.Close()
	// _, err := conn.Do("SET", "token:"+account, token, "EX", timeout)
	// return err
	ret := db.rd.Set("token:"+account, token, time.Duration(timeout)*time.Second)
	return ret.Err()
}

func (db *DBManager) GetLoginToken(account string) (string, error) {
	// conn := db.rd.Get()
	// defer conn.Close()
	// ret, err := conn.Do("GET", "token:"+account)
	// return redis.String(ret, err)
	ret := db.rd.Get("token:" + account)
	return ret.Result()
}

var pubsub *redis.PubSub

func (db *DBManager) CreatePubSub(usermsgcb func(uint64, []byte), roommsgcb func(uint64, []byte), roomadminmsgcb func(uint64, []byte), appmsgcb func(string, []byte), appzonemsgcb func(string, string, []byte)) <-chan *redis.Message {
	pubsub = db.rd.Subscribe()
	msgchan := pubsub.Channel()
	go msgLoop(msgchan, usermsgcb, roommsgcb, roomadminmsgcb, appmsgcb, appzonemsgcb)
	return msgchan
}

func (db *DBManager) StopPubSub() error {
	return pubsub.Close()
}

func msgLoop(msgchan <-chan *redis.Message, usermsgcb func(uint64, []byte), roommsgcb func(uint64, []byte), roomadminmsgcb func(uint64, []byte), appmsgcb func(string, []byte), appzonemsgcb func(string, string, []byte)) {
	for msg := range msgchan {
		keyarr := strings.Split(msg.Channel, ":")
		count := len(keyarr)
		switch keyarr[0] {
		case "user":
			usermsgcb(Uint64(keyarr[count-1]), []byte(msg.Payload))
		case "room":
			roommsgcb(Uint64(keyarr[count-1]), []byte(msg.Payload))
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

func (db *DBManager) SubUserMsg(uids ...uint64) error {
	struids := []string{}
	for _, uid := range uids {
		struids = append(struids, "user:msg:"+String(uid))
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

func (db *DBManager) SubRoomMsg(rids ...uint64) error {
	strrids := []string{}
	for _, rid := range rids {
		strrids = append(strrids, "room:msg:"+String(rid))
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
