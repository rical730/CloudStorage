package controllers

import (
	"LKJapp/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strconv"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	//如果用户按下了退出按钮，则重新设置cookie
	isExit := c.Input().Get("exit") == "true" //先从网页获取用户点击了退出这个操作，并赋予isExit
	if isExit {
		c.Ctx.SetCookie("uname", "", -1, "/") //设置为-1是立即删除cookie的效果
		c.Ctx.SetCookie("pwd", "", -1, "/")
		c.Ctx.SetCookie("cid", "", -1, "/")
		c.Redirect("/", 301) //重定向到首页
		return               //不需要再继续执行下面渲染login画面了
	}

	//用户通过在Controller的对应方法中设置相应的模板名称，beego会自动的在viewpath目录下查询该文件并渲染
	c.TplNames = "login.html"
}

func (c *LoginController) Post() {
	//c.Ctx.WriteString(fmt.Sprint(c.Input()))
	//return

	uname := c.Input().Get("uname")
	pwd := c.Input().Get("pwd")
	autoLogin := c.Input().Get("autoLogin") == "on"

	result, err := models.CheckLogin(uname, pwd)
	if result == true && err == nil { //如果登录成功
		//更新用户的登录次数和最近登录时间
		_ = models.UpdateWhenLogin(uname, pwd)
		//查询是否需要自动登录
		maxAge := 0
		if autoLogin {
			maxAge = 1<<31 - 1
		}
		//设置Cookie
		cid := "0" //默认从根目录开始显示文件夹列表
		c.Ctx.SetCookie("uname", uname, maxAge, "/")
		c.Ctx.SetCookie("pwd", pwd, maxAge, "/")
		c.Ctx.SetCookie("cid", cid, 0, "/")

		fmt.Println("********************* Some one login **************************")
		fmt.Println("* uname: ", uname)
		fmt.Println("* pwd", pwd)
		fmt.Println("* cid", cid)
		fmt.Println("***************************  end  *****************************")
	} else {
		beego.Error(err)
		// c.Data["IsUserLoginFail"] = true //给登陆界面传递消息，登录失败
		c.Redirect("/login", 302) //登录失败重新跳转到登录页面,因为还是本页面，所以用302不能用301
		return
	}

	c.Redirect("/", 301) //重定向到首页
	return
}

//根据网页的Cookie验证登录信息是否正确
func checkAccount(ctx *context.Context) bool {
	ck, err := ctx.Request.Cookie("uname")
	if err != nil {
		return false
	}
	uname := ck.Value

	ck, err = ctx.Request.Cookie("pwd")
	if err != nil {
		return false
	}
	pwd := ck.Value
	//验证账号和密码是否正确
	result, _ := models.CheckLogin(uname, pwd)

	return result
}

//根据网页Cookie获取用户名
func getUnameByCookie(ctx *context.Context) (string, error) {
	ck, err := ctx.Request.Cookie("uname")
	if err != nil {
		return "admin", err
	}
	uname := ck.Value
	return uname, err
}

//根据网页Cookie获取当前应该显示哪一级别的文件夹
func getCidByCookie(ctx *context.Context) (int64, error) {
	ck, err := ctx.Request.Cookie("cid")
	if err != nil {
		return 0, err
	}
	cid, err := strconv.ParseInt(ck.Value, 10, 64)
	return cid, err
}

//根据网页请求更改Cookie中的Cid
func updateCidCookie(ctx *context.Context, newcid string) error {
	_, err := ctx.Request.Cookie("cid")
	if err != nil {
		return err
	}

	// ctx.SetCookie("cid", "", -1, "/")//不需要删除也能直接修改
	ctx.SetCookie("cid", newcid, 0, "/")

	return nil
}
