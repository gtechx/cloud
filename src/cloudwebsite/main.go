package main

import (
	_ "cloudwebsite/routers"
	"gtdb"
	"html/template"
	"math/rand"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
	. "github.com/gtechx/base/common"
)

func Add(a, b int) int {
	return a + b
}

func HtmlAttr(attr string) template.HTMLAttr {
	return template.HTMLAttr(attr)
}

func RandString() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return String(r.Uint32())
}

var DBConfig = `{
	"redisaddr":"192.168.93.12:6379",
    "redispwd":"",
    "redisdefaultdb":2,
    "redismaxconn": 10,

    "mysqladdr":"192.168.93.12:3306",
    "mysqluserpwd":"root:ztgame@123",
    "mysqldb":"gtchat",
    "mysqltableprefix":"gtchat_",
    "mysqllogmode":true,
    "mysqlmaxconn":10,

	"DefaultGroupName": "MyFriends"
}`

func main() {
	defer gtdb.Manager().UnInitialize()
	err := gtdb.Manager().Initialize(DBConfig)
	if err != nil {
		println("Initialize DB err:", err.Error())
		return
	}
	// err := gtdb.Manager().InitializeRedis(config.RedisAddr, config.RedisPassword, config.RedisDefaultDB)
	// if err != nil {
	// 	println("InitializeRedis err:", err.Error())
	// 	return
	// }

	// err = gtdb.Manager().InitializeMysql(config.MysqlAddr, config.MysqlUserPassword, config.MysqlDefaultDB, config.MysqlTablePrefix)
	// if err != nil {
	// 	println("InitializeMysql err:", err.Error())
	// 	return
	// }

	beego.AddFuncMap("Add", Add)
	beego.AddFuncMap("HtmlAttr", HtmlAttr)
	beego.AddFuncMap("RandString", RandString)
	beego.Run()
}
