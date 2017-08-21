package controllers

import (
	"github.com/astaxie/beego" //beego包中会初始化一个BeeAPP的应用，初始化一些参数。
	//"strconv"
	//"fmt"
)

//定义Controller，这里定义了一个struct为MainController，充分利用了Go语言的组合的概念，
type HomeController struct {
	beego.Controller //匿名包含了beego.Controller，这样MainController就拥有了beego.Controller的所有方法。
}

/*************************************************************************************************
定义RESTFul方法：
通过匿名组合之后，其实目前的MainController已经拥有了Get、Post、Delete、Put等方法。
这些方法是分别用来对应用户请求的Method函数，如果用户发起的是POST请求，那么就执行Post函数。
所以这里我们定义了MainController的Get方法用来重写继承的Get函数，这样当用户GET请求的时候就会执行该函数。
**************************************************************************************************/
func (c *HomeController) Get() {
	c.Data["IsHome"] = true

	//将login界面验证到的账号密码匹配与否传递给IsLogin这个参数，在T.naavbar.tpl中验证
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.Data["IsAdminLogin"] = admincheckAccount(c.Ctx) //管理员登录验证

	//首页中的官网和联系地址
	c.Data["Website"] = "http://web.pkusz.edu.cn/antl"
	c.Data["Email"] = "kejiao_li@163.com"

	c.TplNames = "home.html"
}
