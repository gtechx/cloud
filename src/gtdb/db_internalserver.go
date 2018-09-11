package gtdb

import (
	//"errors"

	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

//[sorted sets]serverlist pair(count,addr)
//[sets]ttl:addr
var internalServerKeyName string = "internalserver"

//server op
func (db *DBManager) RegisterInternalServer(addr string) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", internalServerKeyName, 0, addr)

	return err
}

func (db *DBManager) UnRegisterInternalServer(addr string) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("ZREM", internalServerKeyName, addr)

	return err
}

func (db *DBManager) IncrByInternalServerClientCount(addr string, count int) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("ZINCRBY", internalServerKeyName, count, addr)

	return err
}

func (db *DBManager) GetInternalServerList() ([]string, error) {
	conn := db.rd.Get()
	defer conn.Close()

	ret, err := conn.Do("ZRANGE", internalServerKeyName, 0, -1)

	if err != nil {
		return nil, err
	}

	return redis.Strings(ret, err)
}

func (db *DBManager) GetInternalServer() (string, error) {
	conn := db.rd.Get()
	defer conn.Close()

	ret, err := conn.Do("ZRANGE", internalServerKeyName, 0, 1)

	if err != nil {
		return "", err
	}

	slist, err := redis.Strings(ret, err)

	if err != nil || len(slist) == 0 {
		return "", err
	}

	return slist[0], nil
}

func (db *DBManager) GetInternalServerCount() (int, error) {
	conn := db.rd.Get()
	defer conn.Close()

	ret, err := conn.Do("ZCARD", internalServerKeyName)

	count, err := redis.Uint64(ret, err)

	return Int(count), err
}

func (db *DBManager) SetInternalServerTTL(addr string, seconds int) error {
	conn := db.rd.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "ttl:"+addr, 0, "EX", seconds)

	return err
}

func (db *DBManager) CheckInternalServerTTL() error {
	return nil
}

func (db *DBManager) VoteInternalServerDie() error {
	return nil
}
