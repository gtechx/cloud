package admin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	. "github.com/gtechx/base/common"
	"gtdb"
)

type AdminBaseController struct {
	beego.Controller
	tbl_admin *gtdb.Admin
}

func (c *AdminBaseController) checkPrivilege() bool {
	// account := String(c.GetSession("account"))
	// tbl_admin, err := gtdb.Manager().GetAdmin(account)

	// if err != nil {
	// 	fmt.Println("error:", err.Error())
	// 	c.Data["error"] = "数据库错误:" + err.Error()
	// 	return false
	// }

	// c.Data["priv"] = tbl_admin

	adminsess := c.GetSession("admin")
	if adminsess != nil {
		tbl_admin, ok := c.GetSession("admin").(*gtdb.Admin)
		if ok {
			c.Data["priv"] = tbl_admin
			c.tbl_admin = tbl_admin
			return true
		}
	}

	return false
}

type AdminController struct {
	AdminBaseController
}

func (c *AdminController) Prepare() {
	account := String(c.GetSession("account"))
	if account == "" || !c.checkPrivilege() || !c.tbl_admin.Adminadmin {
		c.Redirect("/", 302)
		return
	}

	c.Data["account"] = account
	c.Data["isadmin"] = true
	c.Data["nav"] = "adminadmin"
}

func (c *AdminController) Index() {
	c.TplName = "admin/admin.tpl"
}

func (c *AdminController) Create() {
	if c.Ctx.Request.Method == "POST" {
		c.Data["post"] = true
		account := c.GetString("account")
		adminadmin := c.GetString("adminadmin") == "on"
		adminaccount := c.GetString("adminaccount") == "on"
		adminapp := c.GetString("adminapp") == "on"
		adminappdata := c.GetString("adminappdata") == "on"
		adminonline := c.GetString("adminonline") == "on"
		adminmessage := c.GetString("adminmessage") == "on"
		strexpire := c.GetString("expire")
		expire, exerr := time.Parse("01/02/2006", strexpire)

		if strexpire == "" {
			expire = time.Date(2099, 1, 1, 0, 0, 0, 0, time.Local)
		}

		dbManager := gtdb.Manager()
		var tbl_admin *gtdb.Admin
		var flag bool
		var err error

		// if err != nil {
		// 	println(err.Error())
		// 	c.Data["error"] = err.Error()
		// 	goto end
		// }

		if !adminadmin && !adminaccount && !adminapp && !adminappdata && !adminonline && !adminmessage {
			c.Data["error"] = "不能所有权限都为空"
			goto end
		}

		if account == "" {
			c.Data["error"] = "account不能为空"
			goto end
		}

		flag, err = dbManager.IsAccountExists(account)

		if err != nil {
			println(err.Error())
			c.Data["error"] = "数据库错误:" + err.Error()
			goto end
		}

		if !flag {
			c.Data["error"] = "账号 " + account + " 不存在"
			goto end
		}

		flag, err = dbManager.IsAdmin(account)

		if err != nil {
			println(err.Error())
			c.Data["error"] = "数据库错误:" + err.Error()
			goto end
		}

		if flag {
			c.Data["error"] = "账号 " + account + " 已经是管理员"
			goto end
		}

		if strexpire != "" && exerr != nil && expire.Unix() < time.Now().Unix() {
			c.Data["error"] = "管理员权限过期时间不能小于当前时间"
			goto end
		}

		tbl_admin = &gtdb.Admin{Account: account, Adminadmin: adminadmin, Adminaccount: adminaccount, Adminapp: adminapp, Adminappdata: adminappdata, Adminonline: adminonline, Adminmessage: adminmessage, Expire: expire}
		err = dbManager.CreateAdmin(tbl_admin)

		if err != nil {
			println(err.Error())
			c.Data["error"] = "数据库错误"
			goto end
		}

		c.Redirect("index", 302)
		return
	}

end:
	c.TplName = "admin/admin_create.tpl"
}

func (c *AdminController) Update() {
	account := c.GetString("account")
	dbManager := gtdb.Manager()
	var old_admin *gtdb.Admin
	var err error
	if account == "" {
		c.Data["error"] = "account不能为空"
		goto end
	}

	if account == String(c.GetSession("account")) {
		c.Data["error"] = "不能修改自己的权限"
		goto end
	}

	old_admin, err = dbManager.GetAdmin(account)

	if err != nil {
		fmt.Println("error:", err.Error())
		c.Data["error"] = "数据库错误:" + err.Error()
		goto end
	}

	c.Data["admin"] = old_admin
	if c.Ctx.Request.Method == "POST" {
		c.Data["post"] = true
		old_admin.Adminadmin = c.GetString("adminadmin") == "on"
		old_admin.Adminaccount = c.GetString("adminaccount") == "on"
		old_admin.Adminapp = c.GetString("adminapp") == "on"
		old_admin.Adminappdata = c.GetString("adminappdata") == "on"
		old_admin.Adminonline = c.GetString("adminonline") == "on"
		old_admin.Adminmessage = c.GetString("adminmessage") == "on"
		strexpire := c.GetString("expire")
		expire, exerr := time.Parse("01/02/2006", strexpire)

		if !old_admin.Adminadmin && !old_admin.Adminaccount && !old_admin.Adminapp && !old_admin.Adminappdata && !old_admin.Adminonline && !old_admin.Adminmessage {
			c.Data["error"] = "不能所有权限都为空"
			goto end
		}

		// blank_admin := &gtdb.Admin{}

		// new_admin := &gtdb.Admin{Adminadmin: adminadmin, Adminaccount: adminaccount, Adminapp: adminapp, Adminappdata: adminappdata, Adminonline: adminonline, Adminmessage: adminmessage, Expire: expire}

		// if strexpire == "" {
		// 	new_admin.Expire = time.Date(2099, 1, 1, 0, 0, 0, 0, time.Local)
		// }

		// oldt := reflect.TypeOf(*old_admin)
		// oldv := reflect.ValueOf(old_admin).Elem()
		// //newt := reflect.TypeOf(new_account)
		// newv := reflect.ValueOf(new_admin).Elem()
		// //blankt := reflect.TypeOf(old_account)
		// blankv := reflect.ValueOf(blank_admin).Elem()

		// if err != nil {
		// 	println(err.Error())
		// 	c.Data["error"] = err.Error()
		// 	goto end
		// }

		if strexpire != "" && exerr != nil && expire.Unix() < time.Now().Unix() {
			c.Data["error"] = "管理员权限过期时间不能小于当前时间"
			goto end
		}

		// if err != nil {
		// 	fmt.Println("error:", err.Error())
		// 	c.Data["error"] = "数据库错误:" + err.Error()
		// 	goto end
		// }

		// for k := 0; k < oldt.NumField(); k++ {
		// 	//fmt.Printf("%s -- %v \n", t.Filed(k).Name, v.Field(k).Interface())
		// 	if oldv.Field(k).Type().Kind() != reflect.Slice && oldv.Field(k).Interface() != newv.Field(k).Interface() && newv.Field(k).Interface() != blankv.Field(k).Interface() {
		// 		oldv.Field(k).Set(newv.Field(k))
		// 	}
		// }

		if old_admin.Expire.Unix() != expire.Unix() {
			old_admin.Expire = expire
		}

		fmt.Println("updateadmin:", old_admin)
		err = dbManager.UpdateAdmin(old_admin)

		if err != nil {
			fmt.Println("error:", err.Error())
			c.Data["error"] = "数据库错误:" + err.Error()
		}
	}
end:
	c.TplName = "admin/admin_update.tpl"
}

func (c *AdminController) Del() {
	accounts := c.GetStrings("account[]")
	fmt.Println(accounts)
	newaccounts := make([]string, 0)
	for _, account := range accounts {
		if account != String(c.GetSession("account")) {
			newaccounts = append(newaccounts, account)
		}
	}

	errtext := ""

	if len(newaccounts) > 0 {
		err := gtdb.Manager().DelAdmins(newaccounts)
		if err != nil {
			errtext = "数据库错误:" + err.Error()
		}
	}

	c.Ctx.Output.Body([]byte("{\"error\":\"" + errtext + "\"}"))
}

func (c *AdminController) List() {
	index := Int(c.GetString("pageNumber")) - 1 //Int(c.Ctx.Input.Param("0"))
	pagesize := Int(c.GetString("pageSize"))    //Int(c.Ctx.Input.Param("1"))

	dataManager := gtdb.Manager()
	pagenone := "{\"total\":0, \"rows\":[]}"
	adminfilter := &gtdb.AdminFilter{}

	adminfilter.Account = c.GetString("account")
	adminfilter.Adminadmin = c.GetString("adminadmin") == "on"
	adminfilter.Adminaccount = c.GetString("adminaccount") == "on"
	adminfilter.Adminapp = c.GetString("adminapp") == "on"
	adminfilter.Adminappdata = c.GetString("adminappdata") == "on"
	adminfilter.Adminonline = c.GetString("adminonline") == "on"
	adminfilter.Adminmessage = c.GetString("adminmessage") == "on"

	expire, err := time.Parse("01/02/2006", c.GetString("expire"))
	if err == nil {
		adminfilter.Expire = &expire
	}

	totalcount, err := dataManager.GetAdminCount(adminfilter)

	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte(pagenone))
		return
	}

	if totalcount == 0 {
		c.Ctx.Output.Body([]byte(pagenone))
		return
	}

	adminlist, err := dataManager.GetAdminList(index*pagesize, index*pagesize+pagesize-1, adminfilter)

	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte(pagenone))
		return
	}

	pageadmin := PageData{Total: totalcount, Rows: adminlist}
	retjson, err := json.Marshal(pageadmin)
	if err != nil {
		println(err.Error())
		c.Ctx.Output.Body([]byte(pagenone))
		return
	}

	c.Ctx.Output.Body(retjson)
}
