package app

import (
	"github.com/kataras/iris/mvc"
	"github.com/alezh/novel/bootstrap"
	_ "github.com/alezh/novel/app/admin"
	//引入路由
)

var ConfigureMVC = make([]func(*mvc.Application),0)

func init(){

}

func AddRoute(f func(*mvc.Application)){
	ConfigureMVC = append(ConfigureMVC,f)
}

func MvcBind(app *bootstrap.Bootstrapper)  {
	if len(ConfigureMVC)>0{
		mvc.Configure(app, ConfigureMVC...)
	}
}
