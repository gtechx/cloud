package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/astaxie/beego"
	. "github.com/gtechx/base/common"
	"gtdb"
)

type WebAppController struct {
	beego.Controller
}

func (c *WebAppController) Prepare() {
	//c.BaseController.Prepare()
	c.Data["appaccount"] = String(c.GetSession("account"))
}

func (c *WebAppController) Index() {
	count, _ := gtdb.Manager().GetAppCount()
	applist, _ := gtdb.Manager().GetAppList(0, int(count))
	c.Data["applist"] = applist
	c.TplName = "webapp.tpl"
}

func (c *WebAppController) ChatLogin() {
	account := c.GetString("account")
	password := c.GetString("password")
	appname := c.GetString("appname")

	resp, err := http.PostForm("http://127.0.0.1:9001/chatlogin", url.Values{"account": {account}, "password": {password}, "appname": {appname}})
	defer resp.Body.Close()

	if err != nil {
		c.Ctx.Output.Body([]byte("{\"error\":http error " + err.Error() + ", \"errorcode\":1" + "}"))
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.Ctx.Output.Body([]byte("{\"error\":http error " + err.Error() + ", \"errorcode\":1" + "}"))
		} else {
			c.Ctx.Output.Body(body)
		}
	}
}

func (c *WebAppController) ZoneList() {
	appname := c.GetString("appname")
	println("appname:", appname)
	if appname == "" {
		c.Ctx.Output.Body([]byte(""))
		return
	}

	var err error
	var errtext = ""
	var retjson []byte
	var zonelist []*gtdb.AppZone
	var pagezone PageData

	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsAppExists(appname)
	if err != nil {
		errtext = "数据库错误:" + err.Error()
		goto end
	}
	if !flag {
		errtext = "appname " + appname + " not exists!"
		goto end
	}

	zonelist, err = dbMgr.GetAppZoneList(appname)

	if err != nil {
		println(err.Error())
		errtext = "数据库错误:" + err.Error()
		goto end
	}

	pagezone = PageData{Total: uint64(len(zonelist)), Rows: zonelist}
	retjson, err = json.Marshal(pagezone)
	if err != nil {
		println(err.Error())
		errtext = "json解析错误:" + err.Error()
		goto end
	}

	c.Ctx.Output.Body(retjson)
	return
end:
	c.Ctx.Output.Body([]byte("{\"error\":" + errtext + "}"))
}
