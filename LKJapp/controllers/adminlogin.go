package controllers

import (
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type AdminloginController struct {
	beego.Controller
}

func (c *AdminloginController) Get() {
	//如果用户按下了退出按钮，则重新设置cookie
	adminisExit := c.Input().Get("adminexit") == "true" //先从网页获取用户点击了退出这个操作，并赋予isExit
	if adminisExit {
		c.Ctx.SetCookie("adminname", "", -1, "/") //设置为-1是立即删除cookie的效果
		c.Ctx.SetCookie("adminpwd", "", -1, "/")
		c.Redirect("/", 301) //重定向到首页
		return               //不需要再继续执行下面渲染login画面了
	}

	//用户通过在Controller的对应方法中设置相应的模板名称，beego会自动的在viewpath目录下查询该文件并渲染
	c.TplNames = "adminlogin.html"
}

func (c *AdminloginController) Post() {
	//c.Ctx.WriteString(fmt.Sprint(c.Input()))
	//return

	adminname := c.Input().Get("adminname")
	adminpwd := c.Input().Get("adminpwd")
	adminautoLogin := c.Input().Get("adminautoLogin") == "on"

	if beego.AppConfig.String("adminname") == adminname &&
		beego.AppConfig.String("adminpwd") == adminpwd {
		maxAge := 0
		if adminautoLogin {
			maxAge = 1<<31 - 1
		}
		c.Ctx.SetCookie("adminname", adminname, maxAge, "/")
		c.Ctx.SetCookie("adminpwd", adminpwd, maxAge, "/")
	}

	c.Redirect("/", 301) //重定向到首页
	return
}

func admincheckAccount(ctx *context.Context) bool {
	ck, err := ctx.Request.Cookie("adminname")
	if err != nil {
		return false
	}
	adminname := ck.Value

	ck, err = ctx.Request.Cookie("adminpwd")
	if err != nil {
		return false
	}
	adminpwd := ck.Value

	return beego.AppConfig.String("adminname") == adminname &&
		beego.AppConfig.String("adminpwd") == adminpwd
}
