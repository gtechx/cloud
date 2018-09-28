package gtdb

import (
	"github.com/garyburd/redigo/redis"
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

func (db *DBManager) IsPresenceExists(id, from uint64) (bool, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("EXISTS", "user:presence:"+String(id)+":"+String(from))
	return redis.Bool(ret, err)
}

func (db *DBManager) AddPresence(from, to uint64) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("SET", "user:presence:"+String(to)+":"+String(from), "EX", 60*60*24*7) //记录到目的地用户presence列表
	return err
}

func (db *DBManager) RemovePresence(id, from uint64) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", "user:presence:"+String(id)+":"+String(from))
	return err
}

// func (db *DBManager) GetAllPresence(id uint64) (map[string]string, error) {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGETALL", "user:presence:"+String(id))
// 	return redis.StringMap(ret, err) //.ByteSlices(ret, err)
// }

func (db *DBManager) AddRoomPresence(rid, appdataid uint64) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", "room:presence:"+String(rid), appdataid)
	//_, err := conn.Do("SET", "roompresence:"+String(rid)+":"+String(appdataid), "", "ex 259200") //记录到目的地用户presence列表
	return err
}

func (db *DBManager) RemoveRoomPresence(rid, appdataid uint64) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("SREM", "room:presence:"+String(rid), appdataid)
	//_, err := conn.Do("DEL", "roompresence:"+String(rid)+":"+String(appdataid))
	return err
}

func (db *DBManager) GetAllRoomPresence(rid uint64) ([]uint64, error) {
	conn := db.rd.Get()
	defer conn.Close()
	uids := []uint64{}

	ret, err := conn.Do("SMEMBERS", "room:presence:"+String(rid))
	if err != nil {
		return uids, err
	}
	retarr, err := redis.Values(ret, err)
	if err != nil {
		return uids, err
	}
	err = redis.ScanSlice(retarr, &uids)
	return uids, err //.ByteSlices(ret, err)
}

func (db *DBManager) IsRoomPresenceExists(rid, appdataid uint64) (bool, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", "room:presence:"+String(rid), appdataid)
	//ret, err := conn.Do("EXISTS", "roompresence:"+String(rid)+":"+String(appdataid))
	return redis.Bool(ret, err)
}

func (db *DBManager) PullOnlineMessage(serveraddr string) ([]byte, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("LPOP", "message:"+serveraddr)
	return redis.Bytes(ret, err)
}

func (db *DBManager) GetOfflineMessage(id uint64) ([][]byte, error) {
	conn := db.rd.Get()
	defer conn.Close()

	ret, err := conn.Do("LRANGE", "message:offline:"+String(id), 0, -1)
	datalist, err := redis.ByteSlices(ret, err)
	conn.Do("LTRIM", "message:offline:"+String(id), len(datalist), -1)

	return datalist, err
	// if err != nil {
	// 	return nil, err
	// }

	// retarr, err := redis.Values(ret, nil)

	// if err != nil {
	// 	return nil, err
	// }

	// msglist := [][]byte{}
	// for i := 1; i < len(retarr); i++ {
	// 	msglist = append(msglist, Bytes(retarr[i]))
	// }

	// return msglist, err
}

func (db *DBManager) SendMsgToServer(serveraddr string, msg []byte) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", "message:"+serveraddr, msg)
	return err
}

func (db *DBManager) SendMsgToUserOffline(to uint64, data []byte) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", "message:offline:"+String(to), data)
	return err
}

func (db *DBManager) AddUserMsgHistory(to uint64, timestamp int64, data []byte) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", "user:message:history:"+String(to), timestamp, data)
	return err
}

func (db *DBManager) GetUserMsgHistory(to uint64, mintimestamp int64) ([][]byte, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("ZREVRANGEBYSCORE ", "user:message:history:"+String(to), "+inf", mintimestamp, "LIMIT", 0, 100)
	return redis.ByteSlices(ret, err)
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
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", "room:message:history:"+String(rid), timestamp, msg)
	return err
}

func (db *DBManager) GetRoomMsgHistory(rid uint64, mintimestamp int64) ([][]byte, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("ZREVRANGEBYSCORE ", "room:message:history:"+String(rid), "+inf", mintimestamp, "LIMIT", 0, 100)
	return redis.ByteSlices(ret, err)
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
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("SET", "token:"+account, token, "EX", timeout)
	return err
}

func (db *DBManager) GetLoginToken(account string) (string, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("GET", "token:"+account)
	return redis.String(ret, err)
}
