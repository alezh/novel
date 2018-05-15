package app

import (
	"github.com/alezh/novel/bootstrap"
	"time"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	//_ "github.com/alezh/novel/lib"
)

//var WebApp = NewApp()

var webApp  *bootstrap.Bootstrapper


func NewApp() *bootstrap.Bootstrapper{

	webApp = bootstrap.New("Novel", "alezh.liu@gmail.com")
	webApp.SetupViewsNotLayout("./views")
	webApp.SetupSessions(24*time.Hour,
		[]byte("the-big-and-secret-fash-key-here"),
		[]byte("lot-secret-of-characters-big-too"),
		//[]byte("877253BEFAD283010E5F415D828543D1"),
		//[]byte("7D0AD2BD85FAA6413E4FA8E5AE761119"),
	)
	webApp.SetupErrorHandlers()
	//app.Favicon(bootstrap.StaticAssets + bootstrap.Favicon)
	webApp.StaticWeb("/public", "./public")
	webApp.Use(recover.New())
	webApp.Use(logger.New())
	MvcBind()
	return webApp
}
