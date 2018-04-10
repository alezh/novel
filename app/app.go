package app

import (
	"github.com/alezh/novel/bootstrap"
	"time"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/go-xorm/xorm"
	"github.com/alezh/novel/storage"
)

var WebApp = NewApp()

type (

	WebSystem struct {
		App     *bootstrap.Bootstrapper
	}

)

func NewApp() *WebSystem{
	app := bootstrap.New("Network novel system", "alezh.liu@gmail.com")
	app.SetupViewsNotLayout("../github.com/alezh/novel/views")
	app.SetupSessions(24*time.Hour,
		[]byte("the-big-and-secret-fash-key-here"),
		[]byte("lot-secret-of-characters-big-too"),
	)
	app.SetupErrorHandlers()
	app.Favicon(bootstrap.StaticAssets + bootstrap.Favicon)
	app.StaticWeb("/public", bootstrap.StaticAssets)
	app.Use(recover.New())
	app.Use(logger.New())
	//app.Configure(Routes)
	return &WebSystem{app}
}

func (w *WebSystem)BindMvcApp()  {

}