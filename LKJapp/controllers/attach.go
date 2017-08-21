package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"net/url" //里面有反编码处理
	"os"
)

type AttachController struct {
	beego.Controller
}

func (c *AttachController) Get() {
	fmt.Println("***********正在处理附件下载....")
	//QueryUnescape处理中文名文件，RequestURI截取第一个/后面的所有字符串
	filePath, err := url.QueryUnescape(c.Ctx.Request.RequestURI[1:])
	if err != nil {
		fmt.Println("***********Beego Error:获取文件URI失败")
		c.Ctx.WriteString(err.Error())
		return
	}
	fmt.Println("filePath：", filePath)

	f, err := os.Open(filePath) //f是输入流，实现了io.Read接口
	if err != nil {
		fmt.Println("***********Beego Error:打开服务器上文件路径失败")
		c.Ctx.WriteString(err.Error())
		return
	}
	defer f.Close()

	_, err = io.Copy(c.Ctx.ResponseWriter, f) //输出流是HTTP相应c.Ctx.ResponseWriter，f是输入流
	if err != nil {
		fmt.Println("***********Beego Error:获取文件URI失败")
		c.Ctx.WriteString(err.Error())
		return
	}
	fmt.Println("***********文件获取成功")
}

//读取文件大小，用于category.go文件中上传文件之后检查获取文件大小
func readFileSize(infilename string) (int64, error) {
	file, err := os.Open(infilename)
	defer file.Close()
	if err != nil {
		fmt.Println("failed to open:", infilename)
		return 0, err
	}

	//调用file的方法获取文件信息，返回fileStat结构体类型和错误参数
	finfo, err1 := file.Stat()
	defer file.Close()
	if err1 != nil {
		fmt.Println("get file info failed:", file)
		return 0, err1
	}
	size := int64(finfo.Size())
	return size, nil
}
