package app

import (
	"github.com/kataras/iris/mvc"
	"github.com/alezh/novel/app/admin/controller"
	"github.com/alezh/novel/storage"
)

func init()  {

}

func AdminRoutes(m *mvc.Application)  {
	admin := m.Party("/Admin")
	admin.Register(
		storage.Source,
		webApp.Sessions.Start,
	)
	admin.Handle(new(controller.AdminController))
}

//func Routes(b *bootstrap.Bootstrapper) {
//	b.Get("/stsyem", GetIndexHandler)
//	b.Get("/follower/{id:long}", GetFollowerHandler)
//	b.Get("/following/{id:long}", GetFollowingHandler)
//	b.Get("/like/{id:long}", GetLikeHandler)
//}