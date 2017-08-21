package controllers

import (
	"LKJapp/models"
	"fmt"
	"github.com/astaxie/beego" //beego包中会初始化一个BeeAPP的应用，初始化一些参数。
)

type UserController struct {
	beego.Controller
}

func (c *UserController) Get() {
	//获取/user网页传递过来的userop参数
	userop := c.Input().Get("userop")
	switch userop {
	//如果/user网页传递过来的userop参数是添加操作
	case "add":
		username := c.Input().Get("username")
		if len(username) == 0 {
			break
		}
		err := models.AddUser(username)
		if err != nil {
			fmt.Println("***********Beego Error:管理员新增测试用户失败！")
			//beego.Error(err)
		}
		fmt.Println("***********Beego Error:管理员新增测试用户成功！")

		c.Redirect("/user", 302)
		return

	//如果/user网页传递过来的userop参数的删除操作
	case "del":
		userid := c.Input().Get("userid")
		if len(userid) == 0 {
			break
		}
		err := models.DeleteUser(userid)
		if err != nil {
			fmt.Println("***********Beego Error:管理员删除用户失败！")
			//beego.Error(err)
		}
		fmt.Println("***********Beego Error:管理员删除一个用户成功！用户id=", userid)

		c.Redirect("/user", 302) //重定向回去，刷新列表
		return
	}

	c.Data["IsUser"] = true
	//将login界面验证到的账号密码匹配与否传递给IsLogin这个参数，在T.naavbar.tpl中验证
	c.Data["IsAdminLogin"] = admincheckAccount(c.Ctx)
	//用户通过在Controller的对应方法中设置相应的模板名称，beego会自动的在viewpath目录下查询该文件并渲染
	c.TplNames = "user.html"

	//获取数据库中所有的用户
	var err error
	c.Data["Users"], err = models.GetAllUsers()
	if err != nil {
		fmt.Println("***********Beego Error:管理员获取所有用户信息成功！")
		//beego.Error(err)
	}

	//---------------------------------------------------------------------------------
	op := c.Input().Get("op")

	switch op {
	case "add":
		name := c.Input().Get("name")
		if len(name) == 0 { //文件夹名称不能为空
			break
		}

		uname, err := getUnameByCookie(c.Ctx)
		if uname != "admin" {
			if err != nil {
				fmt.Println("***********Beego Error:新增文件夹前获取管理员名字失败！")
				break
			}
		}

		err = models.AddCategory(name, uname, 0) //管理员新增的测试文件夹全部放在根目录下
		if err != nil {
			fmt.Println("**********Beego Error:管理员新增文件夹失败！")
			//beego.Error(err)
		}
		fmt.Println("**********Beego Error:管理员新增一个文件夹！")
		c.Redirect("/user", 302)
		return

	case "del":
		id := c.Input().Get("id")
		if len(id) == 0 {
			break
		}
		err := models.DeleteCategory(id, 0) //管理员删除文件夹不需要cid
		if err != nil {
			fmt.Println("**********Beego Error:管理员删除文件夹失败！")
			//beego.Error(err)
		}
		fmt.Println("**********Beego Error:管理员删除文件夹成功！id=", id)
		c.Redirect("/user", 302) //重定向回去，刷新列表
		return
	}

	//获取数据库中所有的文件夹
	c.Data["Categories"], err = models.AdminGetAllCategories()
	fmt.Println("info:  succeed!")
	if err != nil {
		fmt.Println("**********Beego Error:管理员获取所有文件夹失败！")
		beego.Error(err)
	}
	fmt.Println("**********Beego Error:管理员获取所有文件夹成功！")

	c.Data["Topics"], err = models.AdminGetAllTopics()
	if err != nil {
		fmt.Println("**********Beego Error:用户获取当前文件信息失败！")
		beego.Error(err)
	} else {
		fmt.Println("**********Beego Info:用户获取当前文件信息成功！")
	}
	//---------------------------------------------------------------------------------

}
