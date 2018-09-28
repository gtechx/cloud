package controllers

import (
	"encoding/json"
	"time"

	. "github.com/gtechx/base/common"
	"gtdb"
)

type AppController struct {
	BaseController
}

func (c *AppController) Prepare() {
	c.BaseController.Prepare()
	c.Data["nav"] = "userapp"
}

func (c *AppController) Index() {
	c.TplName = "app.tpl"
}

func (c *AppController) Create() {
	if c.Ctx.Request.Method == "POST" {
		appname := c.GetString("appname")
		desc := c.GetString("desc")
		share := c.GetString("share")

		println("appcreate ", appname, desc, share)
		c.Data["post"] = true

		dataManager := gtdb.Manager()
		var flag bool
		var err error
		var tbl_app *gtdb.App

		if appname == "" {
			c.Data["error"] = "应用名字不能为空"
			goto end
		}

		flag, err = dataManager.IsAppExists(appname)

		if err != nil {
			println(err.Error())
			c.Data["error"] = "数据库错误"
			goto end
		}

		if flag {
			c.Data["error"] = "应用名字已经存在"
			goto end
		}

		if share != "" {
			flag, err = dataManager.IsAppExists(share)

			if err != nil {
				println(err.Error())
				c.Data["error"] = "数据库错误"
				goto end
			}

			if !flag {
				c.Data["error"] = "共享数据应用名字不存在"
				goto end
			}
		}

		tbl_app = &gtdb.App{Appname: appname, Owner: c.account, Desc: desc, Share: share}
		err = dataManager.CreateApp(tbl_app)

		if err != nil {
			println(err.Error())
			c.Data["error"] = "数据库错误"
			goto end
		}

		c.Redirect("index", 302)
		return
	}
end:
	c.TplName = "appcreate.tpl"
}

func (c *AppController) Update() {
	appname := c.GetString("appname")
	dataManager := gtdb.Manager()
	c.Data["appname"] = appname

	if appname == "" {
		c.Data["error"] = "应用名字为空"
		goto end
	}

	if c.Ctx.Request.Method == "POST" {
		desc := c.GetString("desc")
		share := c.GetString("share")
		c.Data["post"] = true

		println(appname, desc, share)
		err := dataManager.SetShareApp(appname, share)
		if err != nil {
			println("dataManager.SetShareApp ", err.Error())
			c.Data["error"] = "数据库错误:" + err.Error()
			goto end
		}

		err = dataManager.SetAppField(appname, "desc", desc)

		if err != nil {
			println("dataManager.SetAppField ", err.Error())
			c.Data["error"] = "设置应用描述时数据库错误:" + err.Error()
			goto end
		}

		c.Data["desc"] = desc
		c.Data["share"] = share
	} else {
		app, err := dataManager.GetApp(appname)

		if err == nil {
			c.Data["desc"] = app.Desc
			c.Data["share"] = app.Share
		} else {
			println(err.Error())
			c.Data["error"] = "数据库错误:" + err.Error()
		}
	}
end:
	c.TplName = "appmodify.tpl"
}

func (c *AppController) Del() {
	appnames := c.GetStrings("appname[]")

	errtext := ""

	if len(appnames) > 0 {
		err := gtdb.Manager().DeleteApps(appnames)
		if err != nil {
			errtext = "数据库错误:" + err.Error()
		}
	}

	c.Ctx.Output.Body([]byte("{\"error\":\"" + errtext + "\"}"))
}

func (c *AppController) List() {
	index := Int(c.GetString("pageNumber")) - 1 //Int(c.Ctx.Input.Param("0"))
	pagesize := Int(c.GetString("pageSize"))    //Int(c.Ctx.Input.Param("1"))

	println("pageNumber:", index, " pageSize:", pagesize)

	appfilter := &gtdb.AppFilter{}
	appfilter.Appname = c.GetString("appnamefilter")
	appfilter.Desc = c.GetString("descfilter")
	appfilter.Share = c.GetString("sharefilter")

	cbdate, err := time.Parse("01/02/2006", c.GetString("createbegindate"))
	if err == nil {
		appfilter.Createbegindate = &cbdate
	}
	cedate, err := time.Parse("01/02/2006", c.GetString("createenddate"))
	if err == nil {
		appfilter.Createenddate = &cedate
	}

	dataManager := gtdb.Manager()
	totalcount, err := dataManager.GetAppCountByAccount(c.account, appfilter)

	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte("[]"))
		return
	}

	applist, err := dataManager.GetAppListByAccount(c.account, index*pagesize, index*pagesize+pagesize-1, appfilter)

	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte("[]"))
		return
	}

	pageapp := PageData{Total: totalcount, Rows: applist}
	retjson, err := json.Marshal(pageapp)
	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte("[]"))
		return
	}

	c.Ctx.Output.Body(retjson)
}
