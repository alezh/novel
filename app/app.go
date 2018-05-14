package app

import (
	"github.com/alezh/novel/bootstrap"
	"time"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	//"github.com/kataras/iris/mvc"
)

//var WebApp = NewApp()

var webApp  *bootstrap.Bootstrapper


func NewApp() *bootstrap.Bootstrapper{
	webApp = bootstrap.New("Novel", "alezh.liu@gmail.com")
	webApp.SetupViewsNotLayout("./views")
	webApp.SetupSessions(24*time.Hour,
		[]byte("the-big-and-secret-fash-key-here"),
		[]byte("lot-secret-of-characters-big-too"),
	)
	webApp.SetupErrorHandlers()
	//app.Favicon(bootstrap.StaticAssets + bootstrap.Favicon)
	webApp.StaticWeb("/public", "./public")
	webApp.Use(recover.New())
	webApp.Use(logger.New())
	MvcBind()
	return webApp
}
