package gtdb

import (
	"encoding/json"
	"fmt"
	"time" //"github.com/garyburd/redigo/redis"

	"github.com/go-redis/redis"
	. "github.com/gtechx/base/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type config struct {
	RedisAddr      string `json:"redisaddr"`
	RedisPassword  string `json:"redispwd"`
	RedisDefaultDB uint64 `json:"redisdefaultdb"`
	RedisMaxConn   int    `json:"redismaxconn"`

	MysqlAddr         string `json:"mysqladdr"`
	MysqlUserPassword string `json:"mysqluserpwd"`
	MysqlDB           string `json:"mysqldb"`
	MysqlTablePrefix  string `json:"mysqltableprefix"`
	MysqlLogMode      bool   `json:"mysqllogmode"`
	MysqlMaxConn      int    `json:"mysqlmaxconn"`

	DefaultGroupName string `json:"defaultgroupname"`
}

type Redis struct {
	*redis.Client

	serverAddr     string
	serverPassword string
	defaultDB      uint64
}

func (rdm *Redis) Initialize(saddr, spass string, defaultdb uint64) error {
	rdm.serverAddr = saddr
	rdm.serverPassword = spass
	rdm.defaultDB = defaultdb

	rdm.Client = redis.NewClient(&redis.Options{
		Addr:     saddr,
		Password: spass,          // no password set
		DB:       int(defaultdb), // use default DB
	})

	return nil
}

func (rdm *Redis) UnInitialize() error {
	var err error
	if rdm.Client != nil {
		err = rdm.Client.Close()
	}
	return err
}

// func (rdm *Redis) redisDial() (redis.Conn, error) {
// 	c, err := redis.Dial("tcp", rdm.serverAddr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if rdm.serverPassword != "" {
// 		if _, err := c.Do("AUTH", rdm.serverPassword); err != nil {
// 			c.Close()
// 			return nil, err
// 		}
// 	}
// 	if _, err := c.Do("SELECT", rdm.defaultDB); err != nil {
// 		c.Close()
// 		return nil, err
// 	}
// 	return c, nil
// }

// func (rdm *Redis) redisOnBorrow(c redis.Conn, t time.Time) error {
// 	if time.Since(t) < time.Minute {
// 		return nil
// 	}
// 	_, err := c.Do("PING")
// 	return err
// }

type Mysql struct {
	*gorm.DB

	serverAddr     string
	serverPassword string
	defaultDB      string
	prefix         string
}

func (mdm *Mysql) Initialize(saddr, user_pass, defaultdb, prefix string) error {
	mdm.serverAddr = saddr
	mdm.serverPassword = user_pass
	mdm.defaultDB = defaultdb
	mdm.prefix = prefix

	db, err := gorm.Open("mysql", user_pass+"@tcp("+saddr+")/"+defaultdb+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		return err
	}

	gorm.DefaultTableNameHandler = mdm.DefaultTableNameHandler

	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)
	db.SingularTable(true) // 全局禁用表名复数

	mdm.DB = db
	return err
}

func (mdm *Mysql) DefaultTableNameHandler(db *gorm.DB, defaultTableName string) string {
	return mdm.prefix + defaultTableName
}

func (mdm *Mysql) UnInitialize() error {
	var err error
	if mdm.DB != nil {
		err = mdm.DB.Close()
	}
	return err
}

type DBManager struct {
	rd       *Redis
	sql      *Mysql
	dbconfig *config
}

var instance *DBManager

func Manager() *DBManager {
	if instance == nil {
		instance = &DBManager{}
	}
	return instance
}

var tblprefix string

func (db *DBManager) Initialize(configjson string) error {
	dbconfig := &config{}
	err := json.Unmarshal([]byte(configjson), dbconfig)
	if err != nil {
		return err
	}

	db.dbconfig = dbconfig

	db.rd = &Redis{}
	db.sql = &Mysql{}
	err = db.rd.Initialize(dbconfig.RedisAddr, dbconfig.RedisPassword, dbconfig.RedisDefaultDB)
	if err != nil {
		return err
	}

	err = db.sql.Initialize(dbconfig.MysqlAddr, dbconfig.MysqlUserPassword, dbconfig.MysqlDB, dbconfig.MysqlTablePrefix)
	if err != nil {
		db.rd.UnInitialize()
		return err
	}
	tblprefix = db.sql.prefix
	return nil
}

// func (db *DBManager) InitializeRedis(saddr, spass string, defaultdb uint64) error {
// 	db.rd = &Redis{}
// 	return db.rd.Initialize(saddr, spass, defaultdb)
// }

// func (db *DBManager) InitializeMysql(saddr, user_pass, defaultdb, prefix string) error {
// 	db.sql = &Mysql{}
// 	return db.sql.Initialize(saddr, user_pass, defaultdb, prefix)
// }

func (db *DBManager) UnInitialize() error {
	var err error
	if db.rd != nil {
		err = db.rd.UnInitialize()
		db.rd = nil
	}
	if db.sql != nil {
		err = db.sql.UnInitialize()
		db.sql = nil
	}
	return err
}

func (db *DBManager) Install() error {
	// conn := db.rd.Get()
	// defer conn.Close()

	// _, err := conn.Do("FLUSHDB")
	ret := db.rd.FlushDB()
	var err error

	if ret.Err() != nil {
		fmt.Println("FlushDB error")
		return ret.Err()
	}

	tx := db.sql.Begin()
	for _, dbtable := range db_tables {
		if err = tx.DropTableIfExists(dbtable).Error; err != nil {
			fmt.Println("DropTableIfExists error")
			tx.Rollback()
			return err
		}
		//db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").
		if err = tx.Set("gorm:table_options", "ENGINE=NDBCLUSTER").CreateTable(dbtable).Error; err != nil {
			fmt.Println("CreateTable error")
			tx.Rollback()
			return err
		}
	}

	//add admin
	tbl_account := &Account{Account: "admin", Password: "20ff31d485cf6f5cf2d3f5becaa4d0e8", Salt: "gXdqyk", Regip: "127.0.0.1"}
	if err = tx.Create(tbl_account).Error; err != nil {
		tx.Rollback()
		return err
	}

	tbl_admin := &Admin{Account: "admin", Adminadmin: true, Adminaccount: true, Adminapp: true, Adminappdata: true, Adminonline: true, Adminmessage: true, Expire: time.Date(2099, 1, 1, 0, 0, 0, 0, time.Local)}
	if err = tx.Create(tbl_admin).Error; err != nil {
		tx.Rollback()
		return err
	}

	//create test data
	if err = db.CreateTestData(tx); err != nil {
		tx.Rollback()
		return err
	}

	//add admin wyq
	tbl_admin = &Admin{Account: "wyq", Adminadmin: true, Adminaccount: true, Adminapp: true, Adminappdata: true, Adminonline: true, Adminmessage: true, Expire: time.Date(2099, 1, 1, 0, 0, 0, 0, time.Local)}
	if err = tx.Create(tbl_admin).Error; err != nil {
		tx.Rollback()
		return err
	}
	//end

	tx.Commit()
	return err
}

func (db *DBManager) CreateTestData(tx *gorm.DB) error {
	var err error
	tbl_account := &Account{Account: "wyq", Password: "edf06a849c9ec19ea725bd3c6c4ce225", Salt: "p99U86", Regip: "127.0.0.1"}
	if err = tx.Create(tbl_account).Error; err != nil {
		return err
	}

	tbl_app := &App{Appname: "test1", Owner: "wyq", Desc: "ddddd", Share: ""}
	if err = tx.Create(tbl_app).Error; err != nil {
		return err
	}

	tbl_zone := &AppZone{Zonename: "aaa", Owner: "test1"}
	if err = tx.Create(tbl_zone).Error; err != nil {
		return err
	}

	tbl_appdata := &AppData{Appname: "test1", Zonename: "aaa", Account: "wyq", Nickname: "wyqtest", Regip: "127.0.0.1"}
	if err = db.CreateTestAppData(tx, tbl_appdata); err != nil {
		return err
	}

	for i := 2; i < 10; i++ {
		account := "wyq" + String(i)
		tbl_account = &Account{Account: account, Password: "edf06a849c9ec19ea725bd3c6c4ce225", Salt: "p99U86", Regip: "127.0.0.1"}
		if err = tx.Create(tbl_account).Error; err != nil {
			return err
		}

		tbl_appdata = &AppData{Appname: "test1", Zonename: "aaa", Account: account, Nickname: account + "Name", Regip: "127.0.0.1"}
		if err = db.CreateTestAppData(tx, tbl_appdata); err != nil {
			return err
		}
	}

	return nil
}

func (db *DBManager) CreateTestAppData(tx *gorm.DB, tbl_appdata *AppData) error {
	tmpdb := tx.Create(tbl_appdata)
	if err := tmpdb.Error; err != nil {
		return err
	}
	var count uint64
	if err := tx.Model(&AccountApp{}).Where("account = ?", tbl_appdata.Account).Where("appname = ?", tbl_appdata.Appname).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		if err := tx.Create(tbl_appdata.toAccountApp()).Error; err != nil {
			return err
		}
	}

	if err := tx.Model(&AccountZone{}).Where("account = ?", tbl_appdata.Account).Where("appname = ?", tbl_appdata.Appname).Where("zonename = ?", tbl_appdata.Zonename).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		if err := tx.Create(tbl_appdata.toAccountZone()).Error; err != nil {
			return err
		}
	}

	if err := tx.Create(&Group{Groupname: db.dbconfig.DefaultGroupName, Dataid: tbl_appdata.ID}).Error; err != nil {
		return err
	}

	return nil
}
