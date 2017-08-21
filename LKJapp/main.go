/************************************************
Beego生成的APP工程
具体模板查找网址：beego.me、http://my.oschina.net/astaxie/blog/124040
具体Web前端框架使用文档：http://v3.bootcss.com/
具体Go函数与接口查询网址：https://gowalker.org/

调试日志：
4.6：
已经完成首页导航栏设计，首页内容暂定；home.go,home.html增加了hander模板和navbar模板，简化网页设计
用Beego的orm功能实现用户注册信息，创建用户信息数据库，增加model.go
在static目录下增加ccs和javascript控件，方便在html设计中调用
完成了登录界面设计，实现登录功能，记住密码功能，退出功能，增加登陆界面无输入提示功能以及返回功能,增加login.go,login

4.10
增加分类文件夹功能，增加category.go

4.11
完成了首页的设计.发现了可以通过判断是否登录来隐藏导航栏中某一项
开始规划，将之前的登录设置成会员登录，验证登录的时候要使用搜索函数，
而管理员登录则要另外设置一个controller，可以仿照之前的会员登录

4.12
完成了用户管理功能，完成了管理员登录与退出功能
完成注册功能，注册后在用户管理界面会显示

4.13
完善了会员登录验证功能，不过关于账号密码错误的人机交互并不太好

4.16
完成了文件夹管理界面，更新了管理员管理界面，加入了一些图标辅助显示
发现了cookies不能为中文

4.17
将用户登录与退出功能做完整，并增加了cid的Cookie，显示已登录用户的名下的文件夹列表
现在只有登录用户才能看到用户名下的文件夹，而管理员界面可以看到所有用户的所有文件夹
cid默认为0，暂时不可改。
完成打开文件夹功能，返回上一级功能，新建文件增加了父文件夹属性
新增文件夹的时候也同步了其直系所有父文件夹的最后修改时间

4.18
完成用户文件上传，下载，删除功能并不完善
补充了删除文件时对父文件夹大小的修改

4.18下午
完善了删除功能
完善了删除文件夹时对所有子文件夹和子文件的删除
列出了所有应该达到的功能
完善了新增文件功能，如果原本同一文件夹下就有相同名字的文件的话，直接替换，并更新文件信息
完善了所有更新文件夹，更新文件对父文件夹带来的时间和大小更新
另外在新增文件，删除文件，删除文件夹时同步对用户容量使用情况的更新
增加了用户容量使用情况的界面化显示

4.18晚上
完善了管理员管理的所有显示功能

4.19
完善了管理员登录，用户登录的导航栏设计，用户登陆后只显示用户信息。管理员登录后只显示管理员管理界面


作者：李科浇   2015年4月18日完工
********************************************************************/

package main

import (
	"LKJapp/controllers"
	"LKJapp/models"
	_ "LKJapp/routers"
	"github.com/astaxie/beego" //beego采用了Go语言内置的模板引擎，所有模板的语法和Go的一模一样
	"github.com/astaxie/beego/orm"
	"os"
)

func init() {
	models.RegisterDB()
}

/**********************************************************************************************
 * 函数原型: main()
 * 函数功能：主函数入口，注册各种路由（输入：无）输出：无
 * 最后修改：2015/4/17
************************************************************************************************/
func main() {
	//新建一个测试用户
	//utest := models.User{0, lkj, lkj, 0}

	//打印orm的信息打印出来，方便调试
	orm.Debug = true
	//自动建表（数据库名称：default意思是使用默认的数据库，只需要打印一次表，是否打印相关模式）
	orm.RunSyncdb("default", false, true)

	//注册beego路由。Router函数的两个参数函数，第一个是路径，第二个是Controller的指针。
	beego.Router("/", &controllers.HomeController{}) //注册了请求/的时候，调用MainController
	beego.Router("/category", &controllers.CategoryController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/register", &controllers.RegisterController{})

	beego.Router("/user", &controllers.UserController{})
	beego.Router("/adminlogin", &controllers.AdminloginController{})

	// 文件处理
	os.Mkdir("attachment", os.ModePerm)
	beego.Router("/attachment/:all", &controllers.AttachController{})

	beego.Run() //内部监听了8080端口:Go默认情况会监听你本机所有的IP上面的8080端口
	//停止服务的话，请按ctrl+c
}
