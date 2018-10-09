package gtdb

import (
	"time"
	//"errors"

	"github.com/go-redis/redis"
	//. "github.com/gtechx/base/common"
)

//[sorted sets]serverlist pair(count,addr)
//[sets]ttl:addr
var internalServerKeyName string = "internalserver"

//server op
func (db *DBManager) RegisterInternalServer(addr string) error {
	// conn := db.rd.Get()
	// defer conn.Close()
	// _, err := conn.Do("ZADD", internalServerKeyName, 0, addr)

	// return err

	ret := db.rd.ZAdd(internalServerKeyName, redis.Z{Score: 0, Member: addr})
	return ret.Err()
}

func (db *DBManager) UnRegisterInternalServer(addr string) error {
	// conn := db.rd.Get()
	// defer conn.Close()
	// _, err := conn.Do("ZREM", internalServerKeyName, addr)

	// return err
	ret := db.rd.ZRem(internalServerKeyName, addr)
	return ret.Err()
}

func (db *DBManager) IncrByInternalServerClientCount(addr string, count int) error {
	// conn := db.rd.Get()
	// defer conn.Close()
	// _, err := conn.Do("ZINCRBY", internalServerKeyName, count, addr)

	// return err
	ret := db.rd.ZIncrBy(internalServerKeyName, float64(count), addr)
	return ret.Err()
}

func (db *DBManager) GetInternalServerList() ([]string, error) {
	// conn := db.rd.Get()
	// defer conn.Close()

	// ret, err := conn.Do("ZRANGE", internalServerKeyName, 0, -1)

	// if err != nil {
	// 	return nil, err
	// }

	// return redis.Strings(ret, err)
	ret := db.rd.ZRange(internalServerKeyName, 0, -1)
	return ret.Result()
}

func (db *DBManager) GetInternalServer() ([]string, error) {
	// conn := db.rd.Get()
	// defer conn.Close()

	// ret, err := conn.Do("ZRANGE", internalServerKeyName, 0, 1)

	// if err != nil {
	// 	return "", err
	// }

	// slist, err := redis.Strings(ret, err)

	// if err != nil || len(slist) == 0 {
	// 	return "", err
	// }

	// return slist[0], nil

	ret := db.rd.ZRange(internalServerKeyName, 0, 1)
	return ret.Result()
}

func (db *DBManager) GetInternalServerCount() (int64, error) {
	// conn := db.rd.Get()
	// defer conn.Close()

	// ret, err := conn.Do("ZCARD", internalServerKeyName)

	// count, err := redis.Uint64(ret, err)

	// return Int(count), err
	ret := db.rd.ZCard(internalServerKeyName)
	return ret.Result()
}

func (db *DBManager) SetInternalServerTTL(addr string, seconds int) error {
	// conn := db.rd.Get()
	// defer conn.Close()

	// _, err := conn.Do("SET", "ttl:"+addr, 0, "EX", seconds)

	// return err
	ret := db.rd.Set("ttl:"+addr, 0, time.Duration(seconds))
	return ret.Err()
}

func (db *DBManager) CheckInternalServerTTL() error {
	return nil
}

func (db *DBManager) VoteInternalServerDie() error {
	return nil
}
