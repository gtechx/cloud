package gtdb

import (
	//"errors"

	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

//[sorted sets]serverlist pair(count,addr)
//[sets]ttl:addr
var chatServerKeyName string = "chatserverlist"

//server op
func (db *DBManager) RegisterChatServer(addr string) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", chatServerKeyName, 0, addr)

	return err
}

func (db *DBManager) UnRegisterChatServer(addr string) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("ZREM", chatServerKeyName, addr)

	return err
}

func (db *DBManager) IncrByChatServerClientCount(addr string, count int) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("ZINCRBY", chatServerKeyName, count, addr)

	return err
}

func (db *DBManager) GetChatServerList() ([]string, error) {
	conn := db.rd.Get()
	defer conn.Close()

	ret, err := conn.Do("ZRANGE", chatServerKeyName, 0, -1)

	if err != nil {
		return nil, err
	}

	return redis.Strings(ret, err)
}

func (db *DBManager) GetChatServer() (string, error) {
	conn := db.rd.Get()
	defer conn.Close()

	ret, err := conn.Do("ZRANGE", chatServerKeyName, 0, -1)

	if err != nil {
		return "", err
	}

	slist, err := redis.Strings(ret, err)

	if err != nil || len(slist) == 0 {
		return "", err
	}

	return slist[0], nil
}

func (db *DBManager) GetChatServerCount() (int, error) {
	conn := db.rd.Get()
	defer conn.Close()

	ret, err := conn.Do("ZCARD", chatServerKeyName)

	count, err := redis.Uint64(ret, err)

	return Int(count), err
}

func (db *DBManager) InitChatServerTTL(serveraddr string, seconds int) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("SET", "ttl:"+serveraddr, "", "EX", seconds)
	return err
}

func (db *DBManager) UpdateChatServerTTL(serveraddr string, seconds int) error {
	conn := db.rd.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("EXPIRE", "ttl:"+serveraddr, seconds)
	//conn.Send("EXPIRE", "onlineuser:"+serveraddr, seconds)
	_, err := conn.Do("EXEC")
	return err
}

func (db *DBManager) IsChatServerAlive(serveraddr string) (bool, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("EXISTS", "ttl:"+serveraddr)
	return redis.Bool(ret, err)
}

func (db *DBManager) VoteChatServerDie() error {
	return nil
}

func (db *DBManager) SaveChatLoginToken(token string, databytes []byte, timeout int) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("SET", "chattoken:"+token, databytes, "EX", timeout)
	return err
}

func (db *DBManager) GetChatToken(token string) ([]byte, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("GET", "chattoken:"+token)
	return redis.Bytes(ret, err)
}

// func (db *DBManager) UpdateChatServerUserTTL(serveraddr string, seconds int) error {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("EXPIRE", "onlineuser:"+serveraddr, seconds)
// 	return err
// }

// func (db *DBManager) AddOnlineUser(serveraddr string, uid uint64) error {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("SADD", "onlineuser:"+serveraddr, uid)
// 	return err
// }

// func (db *DBManager) RemoveOnlineUser(serveraddr string, uid uint64) error {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("SREM", "onlineuser:"+serveraddr, uid)
// 	return err
// }

// func (db *DBManager) GetAllOnlineUser(serveraddr string) ([]interface{}, error) {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("SMEMBERS", "onlineuser:"+serveraddr)
// 	return redis.Values(ret, err)
// }

func (db *DBManager) SendServerEvent(serveraddr string, data []byte) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", "event:"+serveraddr, data)
	return err
}

func (db *DBManager) PullServerEvent(serveraddr string) ([]byte, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("LPOP", "event:"+serveraddr)
	return redis.Bytes(ret, err)
}
