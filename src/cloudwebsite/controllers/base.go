package controllers

import (
	"github.com/astaxie/beego"
	. "github.com/gtechx/base/common"
	"gtdb"
)

type BaseController struct {
	beego.Controller
	account string
}

func (c *BaseController) setPrivilege() {
	account := String(c.GetSession("account"))
	tbl_admin, err := gtdb.Manager().GetAdmin(account)
	if err == nil {
		c.Data["priv"] = tbl_admin
	}
}

func (c *BaseController) Prepare() {
	account := String(c.GetSession("account"))
	if account == "" {
		c.Redirect("/", 302)
		return
	}
	c.setPrivilege()
	c.Data["account"] = account
	c.account = account
}
