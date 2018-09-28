package controllers

import (
	"encoding/json"
	"time"

	"gtdb"

	. "github.com/gtechx/base/common"
)

type OnlineController struct {
	BaseController
}

func (c *OnlineController) Prepare() {
	c.BaseController.Prepare()
	c.Data["nav"] = "useronline"
}

func (c *OnlineController) Index() {
	count, _ := gtdb.Manager().GetAppCountByAccount(c.account)
	applist, _ := gtdb.Manager().GetAppListByAccount(c.account, 0, int(count))
	c.Data["applist"] = applist
	c.TplName = "online.tpl"
}

func (c *OnlineController) List() {
	index := Int(c.GetString("pageNumber")) - 1 //Int(c.Ctx.Input.Param("0"))
	pagesize := Int(c.GetString("pageSize"))    //Int(c.Ctx.Input.Param("1"))

	id, _ := c.GetUint64("id", 0)
	appname := c.GetString("appname")
	zonename := c.GetString("zonename")
	//account := c.GetString("account")
	// if account != c.account {
	// 	account = c.account
	// }

	appdatafilter := &gtdb.OnlineFilter{}
	appdatafilter.Nickname = c.GetString("nickname")
	appdatafilter.Account = c.GetString("account")
	appdatafilter.Sex = c.GetString("sex")
	appdatafilter.Country = c.GetString("country")
	appdatafilter.Serveraddr = c.GetString("serveraddr")
	appdatafilter.Platform = c.GetString("platform")
	appdatafilter.Isjinyan = Int(c.GetString("jinyan"))

	// appdatafilter.Desc = c.GetString("desc")
	// appdatafilter.Sex = c.GetString("sex")
	// appdatafilter.Country = c.GetString("country")
	// appdatafilter.Regip = c.GetString("regip")
	// appdatafilter.Lastip = c.GetString("lastip")
	// lastloginbegindate := c.GetString("lastloginbegindate")
	// lastloginenddate := c.GetString("lastloginenddate")
	// createbegindate := c.GetString("createbegindate")
	// createenddate := c.GetString("createenddate")

	obdate, err := time.Parse("01/02/2006", c.GetString("onlinebegindate"))
	if err == nil {
		appdatafilter.Onlinebegindate = &obdate
	}
	oedate, err := time.Parse("01/02/2006", c.GetString("onlineenddate"))
	if err == nil {
		appdatafilter.Onlineenddate = &oedate
	}

	println("pageNumber:", index, " pageSize:", pagesize)

	dataManager := gtdb.Manager()
	pagenone := "{\"total\":0, \"rows\":[]}"

	if id != 0 {
		appdata, err := dataManager.GetOnline(id)

		if err != nil {
			println(err.Error())
			c.Ctx.Output.Body([]byte(pagenone))
			return
		}

		pageapp := PageData{Total: 1, Rows: []*gtdb.Online{appdata}}
		retjson, err := json.Marshal(pageapp)
		if err != nil {
			println(err.Error())
			c.Ctx.Output.Body([]byte(pagenone))
			return
		}

		c.Ctx.Output.Body(retjson)
		return
	}

	if appname == "" {
		println("appname must not null")
		c.Ctx.Output.Body([]byte(pagenone))
		return
	}

	totalcount, err := dataManager.GetOnlineCount(appname, zonename, appdatafilter)

	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte(pagenone))
		return
	}

	if totalcount == 0 {
		c.Ctx.Output.Body([]byte(pagenone))
		return
	}

	appdatalist, err := dataManager.GetOnlineList(appname, zonename, index*pagesize, index*pagesize+pagesize-1, appdatafilter)

	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte(pagenone))
		return
	}

	pageapp := PageData{Total: totalcount, Rows: appdatalist}
	retjson, err := json.Marshal(pageapp)
	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte(pagenone))
		return
	}

	c.Ctx.Output.Body(retjson)
}

func (c *OnlineController) ZoneList() {
	appname := c.GetString("appname")
	// account := c.GetString("account")
	// if account != c.account {
	// 	account = c.account
	// }
	owner, err := gtdb.Manager().GetAppOwner(appname)

	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte("[]"))
		return
	}

	if owner != c.account {
		println(c.account, " has no privilege to op ", appname)
		c.Ctx.Output.Body([]byte("[]"))
		return
	}

	zonelist, err := gtdb.Manager().GetAppZoneList(appname)

	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte("[]"))
		return
	}

	pagezone := PageData{Total: uint64(len(zonelist)), Rows: zonelist}
	retjson, err := json.Marshal(pagezone)
	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte("[]"))
		return
	}

	c.Ctx.Output.Body(retjson)
}

func (c *OnlineController) Ban() {
	strappdatas := c.GetStrings("appdata[]")
	appdatas := make([]uint64, len(strappdatas))

	for i, strappdata := range strappdatas {
		appdatas[i] = Uint64(strappdata)
	}

	errtext := ""
	err := gtdb.Manager().BanAppDatas(appdatas, &gtdb.AppDataBaned{Dateline: time.Date(2099, 1, 1, 0, 0, 0, 0, time.Local)})

	if err != nil {
		errtext = "数据库错误:" + err.Error()
	}

	c.Ctx.Output.Body([]byte("{\"error\":\"" + errtext + "\"}"))
}

func (c *OnlineController) Jinyan() {
	strappdatas := c.GetStrings("appdata[]")
	appdatas := make([]uint64, len(strappdatas))

	for i, strappdata := range strappdatas {
		appdatas[i] = Uint64(strappdata)
	}

	errtext := ""
	err := gtdb.Manager().JinyanAppDatas(appdatas, &gtdb.AppDataJinyan{Dateline: time.Date(2099, 1, 1, 0, 0, 0, 0, time.Local)})

	if err != nil {
		errtext = "数据库错误:" + err.Error()
	}

	c.Ctx.Output.Body([]byte("{\"error\":\"" + errtext + "\"}"))
}

func (c *OnlineController) UnJinyan() {
	strappdatas := c.GetStrings("appdata[]")
	appdatas := make([]uint64, len(strappdatas))

	for i, strappdata := range strappdatas {
		appdatas[i] = Uint64(strappdata)
	}

	errtext := ""
	err := gtdb.Manager().UnJinyanAppDatas(appdatas)

	if err != nil {
		errtext = "数据库错误:" + err.Error()
	}

	c.Ctx.Output.Body([]byte("{\"error\":\"" + errtext + "\"}"))
}
