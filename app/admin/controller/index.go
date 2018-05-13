package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris"
	"github.com/alezh/novel/modules"
	"github.com/alezh/novel/modules/user"
	"github.com/alezh/novel/config"
	"github.com/alezh/novel/system"
)

type AdminController struct {
	modules.AuthController
	user.User
}

type formValue func(string) string

//在服务器启动前和控制器注册前调用一次，在这里您可以向该控制器添加依赖项，并且只允许主调用方跳过。
func (c *AdminController) BeforeActivation(b mvc.BeforeActivation) {
	// 绑定依赖
	// form传入数据函数
	b.Dependencies().Add(func(ctx iris.Context) formValue { return ctx.FormValue })
}

// GetIndex handles GET:/Admin/index
func (c *AdminController)GetIndex() {
	id, err := c.Session.GetInt64(config.SessionIDKey)
	if err != nil || id <= 0{

	}
	use := new(user.User)
	c.Source.Mysql.Get(use)
	c.Ctx.Text(use.User)
}

// GetBy handles GET:/user/{id:long},
func (c *AdminController)GetBy(userID int64){

}

func (c *AdminController)GetLogin()  mvc.Result{
	return mvc.View{
		Name:"admin/login.html",
		Data:iris.Map{"config":"1"},
	}
}

func (c *AdminController)Post()  {
	c.Ctx.Text("^……^")
}

func (c *AdminController)PostLogin(form formValue) {
	var (
		username = form("username")
		password = form("password")
	)
}

// 服务器开启 POST :/Admin/engine/start
func (c *AdminController)PostEngineStart()  {
	system.SystemInfo.SetConfig("Mode", config.OFFLINE)
	system.SystemInfo.Start()
}