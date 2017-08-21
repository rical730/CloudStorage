{{define "navbar"}}
<a class="navbar-brand navbar-inverse" href="/"><span class="glyphicon glyphicon-cloud" aria-hidden="true"></span>&nbsp;LKJ云盘</a>
<div>
    <ul class="nav navbar-nav navbar-inverse">
    	<li {{if .IsHome}} class="active"{{end}} ><a href="/">首页</a></li>
        {{if .IsLogin}}<li {{if .IsCategory}} class="active"{{end}} ><a href="/category">我的云盘</a></li>{{end}}

        {{if .IsAdminLogin}}<li {{if .IsUser}} class="active"{{end}} ><a href="/user">用户管理</a></li>{{end}}
        <!-- <li {{if .IsCategory}} class="active"{{end}} >{{if .IsLogin}}<a href="/category">云盘</a>{{end}}</li> -->
        <!-- <li {{if .IsTopic}} class="active"{{end}} ><a href="/topic">文件预览</a></li> -->
    </ul>
</div>

<div class="pull-right">
	<ul class="nav navbar-nav navbar-inverse">
		{{if .IsAdminLogin}}
		<!--如果“退出”这个链接被点击，浏览器中将重新载入本页面url\login，并向新窗口中传递“adminexit=true”这个“参数名=参数值”-->
		<li><a href="/adminlogin?adminexit=true">退出管理员</a></li>
		{{else}}
			{{if .IsLogin}}
			<!--如果“退出”这个链接被点击，浏览器中将重新载入本页面url\login，并向新窗口中传递“exit=true”这个“参数名=参数值”-->
				<li><a href="/login?exit=true">退出</a></li>
			{{else}}
				<li><a href="/adminlogin">管理员登录</a></li>
				<li><a href="/login">会员登录</a></li>
				<li><a href="/register">没有账号？注册一个吧-></a></li>
			{{end}}
		{{end}}

		
	</ul>
</div>
<nav class="navbar navbar-default navbar-fixed-bottom navbar-inverse">
  <div class="container">
    <p>@作者：李科浇</p>
  </div>
</nav>
{{end}}