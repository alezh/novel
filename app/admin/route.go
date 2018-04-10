package admin

import (
	"github.com/alezh/novel/app"
	"github.com/kataras/iris/mvc"
	"github.com/alezh/novel/app/admin/controller"
	"github.com/alezh/novel/storage"
)

func init()  {
	app.AddRoute(Routes)
}

func Routes(m *mvc.Application)  {
	admin := m.Party("/Admin")
	admin.Register(
		storage.Source,
		app.WebApp.App.Sessions.Start,
	)
	admin.Handle(new(controller.AdminController))
}

//func Routes(b *bootstrap.Bootstrapper) {
//	b.Get("/stsyem", GetIndexHandler)
//	b.Get("/follower/{id:long}", GetFollowerHandler)
//	b.Get("/following/{id:long}", GetFollowingHandler)
//	b.Get("/like/{id:long}", GetLikeHandler)
//}
