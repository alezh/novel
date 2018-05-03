package controller

import (
	"github.com/alezh/novel/modules"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris"
	"github.com/alezh/novel/system"
)

var install = mvc.Response{Path: "/Index/install", Code: iris.StatusSeeOther}

type IndexController struct {
	modules.BaseController
}

type formValue func(string) string

func (c *IndexController) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Add(func(ctx iris.Context) formValue { return ctx.FormValue })
}

func (c *IndexController)Get() mvc.Result {

	dbStype := system.SystemInfo.GetConfig("DbStype").(string)

	return mvc.View{
		Name:"index.html",
		Data:iris.Map{"config":dbStype},
	}
}

//初始化环境安装系统
func (c *IndexController)GetInstall()  {
	
}