package controllers

import (
	"LKJapp/models"
	"fmt"
	"github.com/astaxie/beego"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Get() {

	//用户通过在Controller的对应方法中设置相应的模板名称，beego会自动的在viewpath目录下查询该文件并渲染
	c.TplNames = "register.html"
}

func (c *RegisterController) Post() {
	//c.Ctx.WriteString(fmt.Sprint(c.Input()))
	//return

	uname := c.Input().Get("uname")
	pwd := c.Input().Get("pwd")
	pwdagain := c.Input().Get("pwdagain")

	err := models.RegisterUser(uname, pwd)
	if err != nil {
		beego.Error(err)
	}

	fmt.Println("********************* some one registes ***********************")
	fmt.Println("* uname: ", uname)
	fmt.Println("* pwd", pwd)
	fmt.Println("* pwdagain", pwdagain)
	fmt.Println("***************************  end  *****************************")

	c.Redirect("/", 301) //重定向到首页
	return
}
