package controllers

import (
	"github.com/astaxie/beego"
	. "github.com/gtechx/base/common"
	"gtdb"
)

type MainController struct {
	beego.Controller
}

// func (c *MainController) Prepare() {
// 	account := String(c.GetSession("account"))

// 	println("MainController Index account ", account)
// 	if account != "" {
// 		c.Redirect("/user/index", 303)
// 		return
// 	}
// }

func (c *MainController) Index() {
	// c.Data["Website"] = "beego.me"
	// c.Data["Email"] = "astaxie@gmail.com"
	account := String(c.GetSession("account"))

	if account != "" {
		c.Redirect("/user/index", 302)
		return
	}

	c.TplName = "index.tpl"
}

func (c *MainController) Register() {
	if c.Ctx.Request.Method == "POST" {
		account := c.GetString("account")
		password := c.GetString("password")

		c.Data["post"] = true

		flag, err := gtdb.Manager().IsAccountExists(account)

		if err != nil {
			c.Data["error"] = "数据库错误:" + err.Error()
			goto end
		}

		if flag || account == "admin" || account == "root" {
			c.Data["error"] = "账号已经存在"
			goto end
		}

		if password == "" {
			c.Data["error"] = "密码不能为空"
			goto end
		}

		salt := GetSalt(beego.AppConfig.DefaultInt("saltcount", 6))
		md5password := GetSaltedPassword(password, salt)
		println("salt:", salt, "password:", password, "md5password:", md5password)
		tbl_account := &gtdb.Account{Account: account, Password: md5password, Salt: salt, Regip: c.Ctx.Input.IP()}
		err = gtdb.Manager().CreateAccount(tbl_account)

		if err != nil {
			c.Data["error"] = "数据库错误:" + err.Error()
		}
	}
end:
	c.TplName = "register.tpl"
}

func (c *MainController) Login() {
	if c.Ctx.Request.Method == "POST" {
		account := c.GetString("account")
		password := c.GetString("password")

		c.Data["post"] = true

		flag, err := gtdb.Manager().IsAccountExists(account)

		if err != nil {
			c.Data["error"] = "数据库错误"
			goto end
		}

		if !flag {
			c.Data["error"] = "账号不存在！"
			goto end
		}

		// uid, err := gtdb.Manager().GetUIDByAccount(account)

		// if err != nil {
		// 	c.Data["error"] = "数据库错误"
		// 	goto end
		// }

		tbl_account, err := gtdb.Manager().GetAccount(account)

		if err != nil {
			c.Data["error"] = "数据库错误"
			goto end
		}

		md5password := GetSaltedPassword(password, tbl_account.Salt)
		if md5password != tbl_account.Password {
			c.Data["error"] = "密码错误"
			goto end
		}

		println("account ", account, " logined success")

		c.Data["account"] = account

		c.SetSession("account", account)

		flag, _ = gtdb.Manager().IsAdmin(account)

		if flag {
			tbl_admin, err := gtdb.Manager().GetAdmin(account)
			if err == nil {
				c.SetSession("admin", tbl_admin)
			}
		}

		c.Redirect("/user/index", 302)
		return
	}
end:
	c.TplName = "login.tpl"
}

func (c *MainController) Install() {
	key := c.GetString("key")

	if key == "testkey" {
		println("start install...")
		err := gtdb.Manager().Install()
		println("end install...")
		if err == nil {
			c.Ctx.Output.Body([]byte("<html><body>install success!<br/><a href=\"/\">点击进入主页</a></body></html>"))
		} else {
			c.Ctx.Output.Body([]byte("{error:" + err.Error() + "}"))
		}
	} else {
		c.Ctx.Output.Body([]byte("{error:only admin can install the db}"))
	}
}
