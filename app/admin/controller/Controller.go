package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris"
	"fmt"
	"github.com/alezh/novel/modules"
	"github.com/alezh/novel/modules/user"
)

type AdminController struct {
	//mvc.BaseController
	modules.BaseController
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
	use := new(user.User)
	c.Source.Mysql.Get(use)
	fmt.Println(use)
	c.Ctx.Text(use.User)
}

// GetBy handles GET:/user/{id:long},
func (c *AdminController)GetBy(userID int64){

}

func (c *AdminController)Get()  {
	fmt.Println(4454545)
}