package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris"
	."github.com/alezh/novel/modules/basics"
	"github.com/alezh/novel/config"
	"github.com/alezh/novel/app/admin/services"
)

type AdminController struct {
	AuthController
	Service services.UserService
}

type formValue func(string) string


//在服务器启动前和控制器注册前调用一次，在这里您可以向该控制器添加依赖项，并且只允许主调用方跳过。
func (c *AdminController) BeforeActivation(b mvc.BeforeActivation) {
	// 绑定依赖
	// form传入数据函数
	b.Dependencies().Add(func(ctx iris.Context) formValue { return ctx.FormValue })
}
// GetBy handles GET:/Admin,
func (c *AdminController)Get() mvc.Result{
	return c.GetIndex()
}
// GetIndex handles GET:/Admin/index
func (c *AdminController)GetIndex() mvc.Result{
	id:= c.Session.GetString(config.SessionIDKey)
	if id != ""{
		return mvc.View{
			Name:"admin/default.html",
			Data:iris.Map{"User":id},
		}
	}
	return mvc.Response{
		Path: "/Admin/login",
	}
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

func (c *AdminController)PostLogin(form formValue) mvc.Result{
	var (
		username = form("username")
		password = form("password")
	)
	if id ,ok := c.Service.GetByUsernameAndPassword(username,password);ok{
		c.Session.SetImmutable(config.SessionIDKey, id)
		return mvc.Response{
			Path: "/Admin/index",
			Code: iris.StatusSeeOther,
		}
	}else{
		return mvc.Response{
			Path: "/Admin/login",
		}
	}
}

