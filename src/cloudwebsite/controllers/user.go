package controllers

type UserController struct {
	BaseController
}

func (c *UserController) Prepare() {
	c.BaseController.Prepare()
	c.Data["nav"] = "user"
}

func (c *UserController) Index() {
	c.TplName = "user.tpl"
}

func (c *UserController) Logout() {
	//c.DelSession("account")
	//c.DelSession("password")
	c.DelSession("account")

	c.Redirect("/", 302)
}
