package models

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Unknown/com"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	//_ "github.com/go-sql-driver/mysql" // import your used driver
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_DB_NAME        = "data/LKJapp.db" //设置数据库的名称
	_SQLITE3_DRIVER = "sqlite3"        //设置驱动的名称
)

//具有反射功能的字段必须是导出字段，即首字母大写

//创建两个结构  分类 文章
//将结构提交给orm，它就会去进行创建表等等其他操作
type Category struct { //分类结构
	Id      int64     //文件夹标识，只要有int32或者int64的ID，orm就会自动辨识这个为文件
	Cid     int64     //这个文件夹属于哪个文件夹，为0的话说明位于根目录下
	Uname   string    //谁拥有这份文件
	Title   string    `orm:"index"` //文件夹名称
	Size    int64     `orm:"index"` //文件夹大小
	Created time.Time `orm:"index"` //创建时间；使用索引
	Updated time.Time `orm:"index"` //最后修改时间；使用索引
}

//文件
type Topic struct {
	Id        int64     //文件标识
	Cid       int64     //这份文件属于哪个文件夹，为0的话说明位于根目录下
	CidString string    //Cid的字符串类型
	Uname     string    //谁拥有这份文件
	Title     string    `orm:"index"` //文件名称
	Size      int64     `orm:"index"` //文件大小
	Created   time.Time `orm:"index"` //创建时间
	Updated   time.Time `orm:"index"` //最后一次修改时间

}

//用户
type User struct {
	Id            int64     //用户标识
	Name          string    `orm:"index"` //用户名
	Pwd           string    //用户密码
	Capacity      int64     //用户云盘容量
	UseCapacity   int64     //用户已使用云盘容量
	RegisteTime   time.Time `orm:"index"` //注册时间
	LastLoginTime time.Time `orm:"index"` //最后一次登录时间
	LoginCount    int64     //登录次数
}

//--------------------注册一个数据库，数据库的路径是data/LKJapp.db--------------------------------

//会在mian.go里的初始化函数调用
func RegisterDB() {
	if !com.IsExist(_DB_NAME) { //检查目录是否存在
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm) //使用这个函数可以一次性创建1a/s/b/r/e/
		os.Create(_DB_NAME)
	}

	orm.RegisterModel(new(Category), new(Topic), new(User)) //注册一个模型
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DR_Sqlite)      //注册驱动
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
	//至此已经完成了一个orm的创建
}

//--------------------为user.go而设置的函数--------------------------------------------------------
func AddUser(name string) error {
	o := orm.NewOrm() //获取一个orm对象
	auser := &User{
		Name:          name,
		RegisteTime:   time.Now(),
		LastLoginTime: time.Now(),
	} //创建一个用户对象，用户名是传进来的参数

	// 查询用户名是否已经存在
	qs := o.QueryTable("user")
	err := qs.Filter("name", name).One(auser) //使用One来获取一个对象而不是一个slice
	if err == nil {
		return err //表示查询匹配成功，此名称已经存在在数据库中，用户不能使用该名称
	}

	// 插入数据
	_, err = o.Insert(auser)
	if err != nil {
		//beego.Error(err)
		return err //表示插入不成功
	}

	return nil
}

func DeleteUser(id string) error {
	uid, err := strconv.ParseInt(id, 10, 64) //10:十进制，64：int64
	if err != nil {
		return err
	}

	o := orm.NewOrm() //获取一个orm对象

	auser := &User{Id: uid}
	_, err = o.Delete(auser)
	return err
}

func GetAllUsers() ([]*User, error) {
	o := orm.NewOrm() //获取一个orm对象

	users := make([]*User, 0)

	qs := o.QueryTable("user")
	_, err := qs.All(&users)
	return users, err
}

//--------------------为register.go而设置的函数--------------------------------------------------------

func RegisterUser(uname, upwd string) error {
	o := orm.NewOrm() //获取一个orm对象
	var uCap int64
	uCap = 2147483648 //给用户的初始容量是2GB
	auser := &User{
		Name:          uname,
		Pwd:           upwd,
		Capacity:      uCap,
		RegisteTime:   time.Now(),
		LastLoginTime: time.Now(),
	} //创建一个用户对象，用户名是传进来的参数

	// 查询用户名是否已经存在
	qs := o.QueryTable("user")
	err := qs.Filter("name", uname).One(auser) //使用One来获取一个对象而不是一个slice
	if err == nil {
		return err //表示查询匹配成功，此名称已经存在在数据库中，用户不能使用该名称
	}

	// 插入数据
	_, err = o.Insert(auser)
	if err != nil {

		fmt.Println("************Beego Error: Register a user fail  !")
		beego.Error(err)
		return err //表示插入不成功
	}

	return nil
}

//--------------------为login.go而设置的函数--------------------------------------------------------
//http://golanghome.com/post/251
//查询账号密码是否正确
func CheckLogin(uname, upwd string) (bool, error) {
	var users []User //创建一个用户对象，用户名是传进来的参数

	//Filter方法实际是将一个或多个条件按照AND的关系来连接起来
	//我这里直接使用orm.Condition的对象来连接字符串。
	var cond *orm.Condition
	cond = orm.NewCondition()
	cond = cond.And("Name", uname)
	cond = cond.And("Pwd", upwd)
	var qs orm.QuerySeter
	qs = orm.NewOrm().QueryTable("user").SetCond(cond)
	cnt, err := qs.All(&users)

	if err == nil {
		if cnt >= 1 { //如果有1个匹配的话说明验证成功，一般只有一个
			return true, err
		}
	}
	return false, err

}

func UpdateWhenLogin(uname, upwd string) error {
	o := orm.NewOrm() //获取一个orm对象
	auser := new(User)
	qs := o.QueryTable("user")
	err := qs.Filter("name", uname).Filter("pwd", upwd).One(auser) //使用One来获取一个对象而不是一个slice
	if err == nil {
		//表示查询匹配成功,更新数据
		auser.LastLoginTime = time.Now()
		auser.LoginCount++
		_, err = o.Update(auser)
	}
	return err
}

//--------------------为Category.go而设置的函数--------------------------------------------------------
func AddCategory(name, uname string, cid int64) error {
	o := orm.NewOrm()
	//一开始为了调试方便全部cid=0,现在已经增加父文件夹属性
	cate := &Category{
		Uname: uname, //谁拥有这份文件
		Cid:   cid,   //这个文件夹属于哪个文件夹，为0的话说明位于根目录下
		Title: name,  //文件夹名称
		//Size: //文件夹大小
		Created: time.Now(),
		Updated: time.Now(),
	}

	// 查询数据
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).Filter("cid", cid).One(cate) //同一文件夹下的文件夹名称不能相同
	if err == nil {
		fmt.Println("***********Beego Error:新增文件夹名已存在！")
		return err
	}

	//对其父文件夹的最后修改时间进行更新
	err = updateFatherCategory(cid, 0)
	if err != nil {
		fmt.Println("***********Beego Error:AddCategory对其父文件夹的信息更新失败！")
		return err
	}
	fmt.Println("***********Beego Error:AddCategory对其父文件夹的信息更新成功！")

	// 插入数据
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}

	return nil
}

//删除用户指定文件和文件夹
func DeleteCategory(id string, cid int64) error { //删除的时候还需要把文件夹下的文件全部删除！！！！！！！！！！！！！！！！！！！！！！
	cateid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println("***********Beego Error:DeleteCategory的id字符串转换成int64失败！")
		return err
	}

	cate := new(Category)
	o := orm.NewOrm()

	// 查询数据获取文件大小
	qs := o.QueryTable("category")
	err = qs.Filter("id", cateid).One(cate)
	if err != nil {
		return err
	}

	//对其父文件夹的最后修改时间进行更新
	err = updateFatherCategory(cid, 0-cate.Size)
	if err != nil {
		fmt.Println("***********Beego Error:DeleteCategory对其父文件夹的信息更新失败！")
		return err
	}
	fmt.Println("***********Beego Info:DeleteCategory对其父文件夹的信息更新成功！")

	//对用户信息进行更新
	err = updateUserCapacity(cate.Uname, 0-cate.Size)
	if err != nil {
		fmt.Println("***********Beego Error:DeleteCategory对用户信息进行更新更新失败！")
		return err
	}
	fmt.Println("***********Beego Info:DeleteCategory对用户信息进行更新更新成功！")

	//先删除文件夹下所有文件和文件夹
	err = deleteSonCateTopic(cateid)
	//最后再删除本身这个文件夹
	_, err = o.Delete(cate)
	return err
}

//给管理员获取所有用户的所有文件夹
func AdminGetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()

	cates := make([]*Category, 0)

	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

//给用户获取只属于他的当前界面文件夹
func GetAllCategories(uname string, cid int64) ([]*Category, error) {
	o := orm.NewOrm()

	cates := make([]*Category, 0)

	qs := o.QueryTable("category")
	_, err := qs.Filter("uname", uname).Filter("cid", cid).All(&cates)
	return cates, err
}

//给用户获取当前界面的文件夹名称currentDir,err:=models.getCtitleByCid(c.Ctx)
func GetCategoryByCid(cid int64) (*Category, error) {
	o := orm.NewOrm()

	cate := new(Category)

	qs := o.QueryTable("category")
	err := qs.Filter("id", cid).One(cate)

	return cate, err
}

//---------------------------为Topic信息而设置的函数--------------
func AddTopic(name, uname string, cid, size int64) error {
	cidstring := strconv.FormatInt(cid, 10)
	o := orm.NewOrm()
	//一开始为了调试方便全部cid=0,现在已经增加父文件夹属性
	topic := &Topic{
		Uname:     uname, //谁拥有这份文件
		Cid:       cid,   //这个文件夹属于哪个文件夹，为0的话说明位于根目录下
		CidString: cidstring,
		Title:     name, //文件夹名称
		Size:      size, //文件大小
		Created:   time.Now(),
		Updated:   time.Now(),
	}

	// 查询数据
	qs := o.QueryTable("topic")
	err := qs.Filter("title", name).Filter("cid", cid).One(topic) //同一文件夹下的文件名称不能相同
	if err == nil {
		fmt.Println("***********Beego Error:新增文件名已存在！直接覆盖原文件！")

		//对其父文件夹的最后修改时间进行更新
		err = updateFatherCategory(cid, size-topic.Size)
		if err != nil {
			fmt.Println("***********Beego Error:AddTopic对其父文件夹的信息更新失败！")
			return err
		}
		fmt.Println("***********Beego Error:AddTopic对其父文件夹的信息更新成功！")

		//对用户信息进行更新
		err = updateUserCapacity(uname, size-topic.Size)
		if err != nil {
			fmt.Println("***********Beego Error:DeleteCategory对用户容量进行更新更新失败！")
			return err
		}
		fmt.Println("***********Beego Info:DeleteCategory对用户容量进行更新更新成功！")

		//更新时间
		topic.Updated = time.Now()
		topic.Size = size

		// 插入数据
		_, err = o.Update(topic)
		if err != nil {
			return err
		}

		return err //后面要改成修改文件参数！！！！！！！！！！！！！
	}

	//对其父文件夹的最后修改时间进行更新
	err = updateFatherCategory(cid, size)
	if err != nil {
		fmt.Println("***********Beego Error:AddTopic对其父文件夹的信息更新失败！")
		return err
	}
	fmt.Println("***********Beego Error:AddTopic对其父文件夹的信息更新成功！")

	//对用户信息进行更新
	err = updateUserCapacity(uname, size)
	if err != nil {
		fmt.Println("***********Beego Error:DeleteCategory对用户容量进行更新更新失败！")
		return err
	}
	fmt.Println("***********Beego Info:DeleteCategory对用户容量进行更新更新成功！")

	// 插入数据
	_, err = o.Insert(topic)
	if err != nil {
		return err
	}

	return nil
}

//删除用户指定文件
func DeleteTopic(id string, cid int64) error {
	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	topic := new(Topic)
	o := orm.NewOrm()

	// 查询数据获取文件大小
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tid).One(topic)
	if err != nil {
		return err
	}

	//对其父文件夹的最后修改时间进行更新
	err = updateFatherCategory(cid, 0-topic.Size)
	if err != nil {
		fmt.Println("***********Beego Error:DeleteCategory对其父文件夹的信息更新失败！")
		return err
	} else {
		fmt.Println("***********Beego Error:DeleteCategory对其父文件夹的信息更新成功！")
	}

	//对用户信息进行更新
	err = updateUserCapacity(topic.Uname, 0-topic.Size)
	if err != nil {
		fmt.Println("***********Beego Error:DeleteCategory对用户容量进行更新更新失败！")
		return err
	}
	fmt.Println("***********Beego Info:DeleteCategory对用户容量进行更新更新成功！")

	// topic := &Topic{Id: tid}

	//先删除文件
	//因为文件命名规则是[LKJyunpan.uname.cname]name
	FileName := "[" + "LKJyunpan" + "." + topic.Uname + "." + topic.CidString + "]" + topic.Title
	err = os.Remove(path.Join("attachment", FileName))
	if err != nil {
		fmt.Println("***********Beego Error:删除服务器上的文件失败！")
		return err
	} else {
		fmt.Println("***********Beego Info:删除服务器上的文件成功，文件名=", FileName)
	}

	//再删除元数据
	_, err = o.Delete(topic)
	return err
}

//给管理员获取所有用户的所有文件
func AdminGetAllTopics() ([]*Topic, error) {
	o := orm.NewOrm()

	topics := make([]*Topic, 0)

	qs := o.QueryTable("topic")
	_, err := qs.All(&topics)
	return topics, err
}

//给用户获取只属于他的当前界面文件
func GetAllTopics(uname string, cid int64) ([]*Topic, error) {
	o := orm.NewOrm()

	topics := make([]*Topic, 0)

	qs := o.QueryTable("topic")
	_, err := qs.Filter("uname", uname).Filter("cid", cid).All(&topics)
	return topics, err
}

//用于func (c *CategoryController) Get()，功能是查看用户容量使用情况，获取比例
func GetCapByUname(uname string) string {
	o := orm.NewOrm()
	auser := new(User)
	qs := o.QueryTable("user")
	err := qs.Filter("name", uname).One(auser)
	if err != nil {
		fmt.Println("***********Beego Error:查看容量的时候查询用户信息失败，用户：", uname)
		return "0"
	}
	pro := int64(float32(auser.UseCapacity) * 100 / float32(auser.Capacity))
	proString := strconv.FormatInt(pro, 10)
	return proString
}

//----------------------------------------------------------------------------------
//----------------------------------------------------------------------------------
//----------------------------------------------------------------------------------
//----------------------------------------------------------------------------------
//----------------------------------------------------------------------------------
//----------------------------------------------------------------------------------

//----------------------------------------本地函数-----------------------------------
//用于AddCategory，功能是更新其所有父文件夹的大小和时间属性
//使用递归算法
func updateFatherCategory(cid, sizeAdd int64) error {
	//如果已经是根目录了，那就没有父文件夹，无需修改
	if cid == 0 {
		return nil
	}

	//否则先修改第一层父文件夹，再往上递归
	o := orm.NewOrm()
	cate := new(Category)
	qs := o.QueryTable("category")
	err := qs.Filter("id", cid).One(cate)
	if err == nil {
		//表示查询匹配成功,更新数据
		cate.Updated = time.Now()
		if (cate.Size + sizeAdd) >= 0 {
			cate.Size += sizeAdd
		}
		_, err = o.Update(cate)
		if err == nil {
			err = updateFatherCategory(cate.Cid, sizeAdd)
		}
	}
	return err
}

//用于DeleteCategory，功能是删除其下所有子子孙孙的文件夹和文件
func deleteSonCateTopic(id int64) error {
	o := orm.NewOrm()
	var err error
	var cntCate, cntTopic int64
	//先把扫描出来的文件全部删除
	topics := make([]*Topic, 0)
	qs1 := o.QueryTable("topic")
	cntTopic, err = qs1.Filter("cid", id).All(&topics)
	if err != nil {
		fmt.Println("***********Beego Error:查询失败，文件id=", id)
		return err
	}
	//如果有文件就删除文件，没有的话可以跳过
	if cntTopic >= 1 {
		for _, Tv := range topics {
			//因为文件命名规则是[LKJyunpan.uname.cname]name
			FileName := "[" + "LKJyunpan" + "." + Tv.Uname + "." + Tv.CidString + "]" + Tv.Title
			//删除文件
			err := os.Remove(path.Join("attachment", FileName))
			if err != nil {
				fmt.Println("***********服务器找不到文件cid=", Tv.Cid, Tv.Title)
			}
			//删除元数据
			// o := orm.NewOrm()
			// qs1 := o.QueryTable("topic")
			// err := qs.Filter("title", name).Filter("cid", cid).One(cate)
			_, err = o.Delete(Tv)
			fmt.Println("***********删除文件cid=", Tv.Cid, Tv.Title)
		}
	}

	//然后再继续往下遍历文件夹
	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	cntCate, err = qs.Filter("cid", id).All(&cates)
	if err != nil {
		fmt.Println("***********Beego Error:查询失败，文件夹id=", id)
		return err
	}
	//如果有文件夹就继续遍历文件夹，没有的话可以返回删除文件夹
	if cntCate >= 1 {
		for _, Cv := range cates {
			//往下搜索文件夹
			err := deleteSonCateTopic(Cv.Id)
			if err != nil {
				return err //如果过程中出错，直接返回
			}
			// o := orm.NewOrm()
			// qs := o.QueryTable("category")
			_, err = o.Delete(Cv)
		}
	}
	return err

	//一直到这步，所有文件都已经删除，删除后从最底层的文件夹开始自动销毁

}

//用于DeleteCategory和DeleteTopic和AddTopic，功能是更新用户容量使用情况
func updateUserCapacity(uname string, size int64) error {
	o := orm.NewOrm()
	auser := new(User)
	qs := o.QueryTable("user")
	err := qs.Filter("name", uname).One(auser)
	if err != nil {
		fmt.Println("***********Beego Error:更新容量的时候查询用户信息失败，用户：", uname)
		return err
	}
	if (auser.UseCapacity + size) >= 0 {
		auser.UseCapacity += size
	}
	_, err = o.Update(auser)
	fmt.Println("***********Beego Error:更新容量成功，用户：", uname, auser.UseCapacity)
	return err
}
