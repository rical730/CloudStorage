# CloudStorage
A cloud storage web app.

## Intruction
 * 用Go语言实现了一个仿百度云盘的网站。（对没错我把它命名为LKJ云盘了）

## Background
 * 当时是为了实现实验室的对象存储系统的文件系统网页访问端，我干脆做成了一个网盘的形式
 * 因为参与的项目使用的是Go语言作为后台开发，为了实现更好的通信和适配性，采用了Go Web实现前端开发。
 * Go是一种新的语言，一种并发的、带垃圾回收的、快速编译的语言。
 * Go是一种编译型语言，它结合了解释型语言的游刃有余，动态类型语言的开发效率，以及静态类型的安全性。
 * 感兴趣的可以看谢大师(astaxie)的[GitBook](https://www.gitbook.com/book/wizardforcel/build-web-application-with-golang/details)，当时是跟着这个入门的。

## Tech
 * 语言：Go语言
 * 框架：Beego(典型MVC，国产框架，很成熟)
 * 前端样式：Bootstrap
 * 编辑工具：sublime text
 * 架构风格：Restful

## Function
 * 首页请求
 * 用户管理
 	* 普通用户
 		* 注册
 		* login
 		* logout
 	* 管理员
 		* login
 		* logout
 * 文件夹操作
 	* 创建文件夹
 	* 删除文件夹
 	* 打开文件夹
 	* 返回上一级文件夹
 * 文件操作
 	* 上传
 	* 下载
 * 管理员操作
 	* 增删查用户
 	* 增删查文件
 	* 增删查文件夹
 * 容量管理
 	* 查看容量使用情况

## Other
 * 利用Cookies实现自动登陆功能
 * 页面实时显示文件或文件夹大小、最后修改时间

## Screen Shot
![](https://github.com/rical730/CloudStorage/blob/master/screenshot/main.png)
