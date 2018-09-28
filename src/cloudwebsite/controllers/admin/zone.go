package admin

import (
	"encoding/json"
	"fmt"
	"strings"

	. "github.com/gtechx/base/common"
	"gtdb"
)

type ZoneController struct {
	AdminBaseController
}

func (c *ZoneController) Prepare() {
	account := String(c.GetSession("account"))
	if account == "" || !c.checkPrivilege() || !c.tbl_admin.Adminapp {
		c.Redirect("/", 302)
		return
	}
	c.Data["account"] = account
}

func (c *ZoneController) Create() {
	errtext := ""
	var retjson []byte
	//var err error
	var zonelist []*gtdb.AppZone
	var tbl_zone *gtdb.AppZone
	if c.Ctx.Request.Method == "POST" {
		appname := c.GetString("appname")
		zonename := c.GetString("zonename")
		account := c.GetString("account")
		println(appname, " ", zonename)
		if zonename == "" {
			c.Ctx.Output.Body([]byte(""))
			return
		}
		zonenamearr := strings.Split(zonename, ";")

		dataManager := gtdb.Manager()

		appowner, err := dataManager.GetAppOwner(appname)
		if err != nil {
			errtext = "数据库错误:" + err.Error()
			goto end
		}
		if appowner != account {
			errtext = appname + "不属于" + account
			goto end
		}
		for _, zname := range zonenamearr {
			flag, err := dataManager.IsAppZoneExists(appname, zname)

			if err != nil {
				errtext = "zone " + zname + " add redis error"
				continue
			}

			if flag {
				errtext = "zone " + zname + " already exists"
				continue
			}

			tbl_zone = &gtdb.AppZone{Zonename: zname, Owner: appname}
			err = dataManager.AddAppZone(tbl_zone)

			if err != nil {
				errtext = "zone " + zname + " add redis error"
			}
		}

		println(errtext)
		zonelist, err = dataManager.GetAppZoneList(appname)

		retjson, err = json.Marshal(zonelist)
		if err != nil {
			println(err.Error())
			c.Ctx.Output.Body([]byte(""))
			return
		}

		c.Ctx.Output.Body(retjson)
		return
	end:
		println(errtext)
		c.Ctx.Output.Body([]byte("{\"error\":" + errtext + "}"))
	}
}

func (c *ZoneController) Del() {
	appname := c.GetString("appname")
	zonenames := c.GetStrings("zonename[]")
	//account := c.GetString("account")
	dataManager := gtdb.Manager()
	fmt.Println(appname, zonenames)
	errtext := ""
	// appowner, err := dataManager.GetAppOwner(appname)
	// if err != nil {
	// 	errtext = "数据库错误:" + err.Error()
	// 	goto end
	// }
	// if appowner != account {
	// 	errtext = "没有权限操作该应用：" + appname
	// 	goto end
	// }

	if len(zonenames) > 0 {
		err := dataManager.RemoveAppZones(appname, zonenames)
		if err != nil {
			errtext = "数据库错误:" + err.Error()
		}
	}

	fmt.Println("errtext:", errtext)
	//end:
	c.Ctx.Output.Body([]byte("{\"error\":\"" + errtext + "\"}"))
}

func (c *ZoneController) List() {
	appname := c.GetString("appname")
	println(appname)
	if appname == "" {
		c.Ctx.Output.Body([]byte(""))
		return
	}

	var err error
	//var errtext = ""
	var retjson []byte
	var zonelist []*gtdb.AppZone

	dbMgr := gtdb.Manager()
	// appowner, err := dbMgr.GetAppOwner(appname)
	// if err != nil {
	// 	errtext = "数据库错误"
	// 	goto end
	// }
	// if appowner != c.account {
	// 	errtext = "没有权限操作该app"
	// 	goto end
	// }

	zonelist, err = dbMgr.GetAppZoneList(appname)

	retjson, err = json.Marshal(zonelist)
	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte("{\"error\":" + err.Error() + "}"))
		return
	}

	c.Ctx.Output.Body(retjson)
	return
	//end:
	//c.Ctx.Output.Body([]byte("{\"error\":" + errtext + "}"))
}
