package controllers

import (
	"LKJapp/models"
	"fmt"
	"github.com/astaxie/beego" //beego包中会初始化一个BeeAPP的应用，初始化一些参数。
	//"github.com/astaxie/beego/context"
	"path"
	"strconv"
)

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get() {
	//解析表单
	op := c.Input().Get("op")

	switch op {
	case "add":
		name := c.Input().Get("name")
		if len(name) == 0 { //文件夹名称不能为空
			break
		}

		uname, err1 := getUnameByCookie(c.Ctx)
		if err1 != nil {
			fmt.Println("***********Beego Error:在新增文件夹前获取用户名cookies失败，可能用户未登陆")
			break
		}
		cid, err2 := getCidByCookie(c.Ctx)
		if err2 != nil {
			fmt.Println("***********Beego Error:在获取当前文件夹列表前获取Cid失败，用户未登录")
			break
		}
		err := models.AddCategory(name, uname, cid)
		if err != nil {
			fmt.Println("**********Beego Error:用户新增文件夹失败！")
			//beego.Error(err)
		} else {
			fmt.Println("**********Beego Info:用户新增文件夹成功！用户名和文件夹名：", uname, name)
		}

		c.Redirect("/category", 302)
		return

	case "del":
		id := c.Input().Get("id")
		if len(id) == 0 {
			break
		}
		cid, err2 := getCidByCookie(c.Ctx)
		if err2 != nil {
			fmt.Println("***********Beego Error:在获取当前文件夹列表前获取Cid失败，用户未登录")
			break
		}
		err := models.DeleteCategory(id, cid)
		if err != nil {
			fmt.Println("**********Beego Error:用户删除文件夹失败！")
			beego.Error(err)
		} else {
			fmt.Println("**********Beego Info:用户删除文件夹成功！文件夹id=", id)
		}

		c.Redirect("/category", 302) //重定向回去，刷新列表
		return

	case "open":
		id := c.Input().Get("id")
		if len(id) == 0 {
			break
		}
		err := updateCidCookie(c.Ctx, id)
		if err != nil {
			fmt.Println("**********Beego Error:更改Cookie中的cid失败！用户未登录")
			beego.Error(err)
		} else {
			fmt.Println("**********Beego Info:用户成功打开文件夹id=", id)
		}

		c.Redirect("/category", 302) //重定向回去，刷新列表
		return

	case "levelup":
		cid := c.Input().Get("cid")
		if len(cid) == 0 {
			break
		}
		err := updateCidCookie(c.Ctx, cid)
		if err != nil {
			fmt.Println("**********Beego Error:更改Cookie中的cid失败！用户未登录")
			beego.Error(err)
		} else {
			fmt.Println("**********Beego Info:用户成功打开文件夹id=", cid)
		}

		c.Redirect("/category", 302) //重定向回去，刷新列表
		return
	case "deltopic":
		tid := c.Input().Get("tid")
		if len(tid) == 0 {
			break
		}

		cid, err2 := getCidByCookie(c.Ctx)
		if err2 != nil {
			fmt.Println("***********Beego Error:在获取当前文件夹列表前获取Cid失败，用户未登录")
			break
		}

		err := models.DeleteTopic(tid, cid)
		if err != nil {
			fmt.Println("**********Beego Error:用户删除文件夹失败！")
			beego.Error(err)
		} else {
			fmt.Println("**********Beego Info:用户删除文件夹成功！文件夹tid=", tid)
		}

		c.Redirect("/category", 302) //重定向回去，刷新列表
		return
	}

	c.Data["IsCategory"] = true
	//将login界面验证到的账号密码匹配与否传递给IsLogin这个参数，在T.naavbar.tpl中验证
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	//用户通过在Controller的对应方法中设置相应的模板名称，beego会自动的在viewpath目录下查询该文件并渲染
	c.TplNames = "category.html"

	//获取数据库中所有该用户对象现在cid的文件夹和文件
	var err error
	var uname string
	var cid int64
	uname, err = getUnameByCookie(c.Ctx)
	if err != nil {
		fmt.Println("***********Beego Error:在获取当前文件夹列表前获取用户名失败，用户未登录")
		return
	}

	cid, err = getCidByCookie(c.Ctx)
	if err != nil {
		fmt.Println("***********Beego Error:在获取当前文件夹列表前获取Cid失败，用户未登录")
		return
	}

	c.Data["Categories"], err = models.GetAllCategories(uname, cid)
	if err != nil {
		fmt.Println("**********Beego Error:用户获取当前文件夹列表失败！")
		beego.Error(err)
	} else {
		fmt.Println("**********Beego Info:用户获取当前列表文件夹成功！")
	}

	c.Data["Topics"], err = models.GetAllTopics(uname, cid)
	if err != nil {
		fmt.Println("**********Beego Error:用户获取当前文件信息失败！")
		beego.Error(err)
	} else {
		fmt.Println("**********Beego Info:用户获取当前文件信息成功！")
	}

	//字段赋值
	c.Data["IsRootDir"] = (cid == 0)
	if cid != 0 {
		currentCate, err := models.GetCategoryByCid(cid)
		if err != nil {
			fmt.Println("**********Beego Error:用户获取当前文件夹名称失败！")
			beego.Error(err)
		}
		c.Data["CurrentCate"] = currentCate
	}
	currentCname := strconv.FormatInt(cid, 10)
	c.Data["CurrentCname"] = currentCname
	c.Data["Uname"] = uname

	capProportion := models.GetCapByUname(uname)
	c.Data["CapProportion"] = capProportion
	fmt.Println("***********哈哈哈哈用户容量是：", capProportion, "%")

}

func (c *CategoryController) Post() {
	// 解析表单

	//上传文件——SaveToFile相比较GetFile，SaveToFile有一个比较好的封装可以直接把文件保存到服务器端而不用多余的代码
	//GetFile可以用来判断是否有文件上传，有的话再调用SaveToFile保存文件
	//文件上传后一般是放在服务器系统的内存里，默认缓存64M，可调整
	_, fh, err := c.GetFile("attachment")
	if err != nil {
		fmt.Println("**********Beego Error:用户选择上传文件失败！")
		beego.Error(err)
	}

	var attachment string
	if fh != nil {
		// 保存附件
		attachment = fh.Filename //string类型
		fmt.Println("**********Beego Info:文件名：")
		beego.Info(attachment)

		//获取当前页面用户信息
		var err error
		var uname, cname string
		var cid, size int64
		uname, err = getUnameByCookie(c.Ctx)
		if err != nil {
			fmt.Println("***********Beego Error:在获取上传文件前获取用户名失败，用户未登录")
			return
		}
		cid, err = getCidByCookie(c.Ctx)
		if err != nil {
			fmt.Println("***********Beego Error:在上传文件前获取Cid失败，用户未登录")
			return
		}
		cname = strconv.FormatInt(cid, 10)

		//在源文件当前目录下新建一个目录，命名为[LKJyunpan.uname.cname]name
		//在删除文件、删除文件夹、下载文件都遵守这个命名规则
		newFileName := "[" + "LKJyunpan" + "." + uname + "." + cname + "]" + attachment

		err = c.SaveToFile("attachment", path.Join("attachment", newFileName)) //这里是client服务器保存文件法则
		if err != nil {
			fmt.Println("**********Beego Error:用户保存上传文件失败！")
			beego.Error(err)
		} else {
			fmt.Println("**********Beego Info:用户保存上传文件成功！")
		}

		size, err = readFileSize(path.Join("attachment", newFileName)) //函数放在attach.go中
		if err != nil {
			fmt.Println("**********Beego Error:获取文件大小失败！")
		}
		fmt.Println("**********Beego Info:文件大小=", size)

		//新增一个Topic记录文件信息
		err = models.AddTopic(attachment, uname, cid, size) //暂时文件大小赋值成

		//err = models.ModifyTopic(attachment, uname, cid)

		if err != nil {
			beego.Error(err)
			fmt.Println("**********Beego Error:用户新增Topic信息失败！")
		} else {
			fmt.Println("**********Beego Info:用户新增文件信息成功！")
		}

	}

	c.Redirect("/category", 302)
}
