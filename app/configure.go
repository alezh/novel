package app

import (
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris"
)

var ConfigureMVC = make([]func(*mvc.Application),0)

func init() {
	AddRoute(AdminRoutes,IndexRoutes)
}


func AddRoute(f... func(*mvc.Application)){
	if len(f)>0{
		ConfigureMVC = append(ConfigureMVC,f...)
	}
}

func MvcBind()  {
	if len(ConfigureMVC)>0{
		mvc.Configure(webApp, ConfigureMVC...)
	}
}

func Configure(app *iris.Application) {
	app.Configure(
		iris.WithoutServerError(iris.ErrServerClosed),
	)
}