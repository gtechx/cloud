package gtdb

import (
	//"errors"

	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

//[sorted sets]serverlist pair(count,addr)
//[sets]ttl:addr
var chatServerKeyName string = "chatserver"

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

func (db *DBManager) SetChatServerTTL(addr string, seconds int) error {
	conn := db.rd.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "ttl:"+addr, 0, "EX", seconds)

	return err
}

func (db *DBManager) CheckChatServerTTL() error {
	return nil
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
