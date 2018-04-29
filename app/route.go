package app

import (
	"github.com/kataras/iris/mvc"
	adminController "github.com/alezh/novel/app/admin/controller"
	indexController "github.com/alezh/novel/app/index/controller"
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
	admin.Handle(new(adminController.AdminController))
}

//func Routes(b *bootstrap.Bootstrapper) {
//	b.Get("/stsyem", GetIndexHandler)
//	b.Get("/follower/{id:long}", GetFollowerHandler)
//	b.Get("/following/{id:long}", GetFollowingHandler)
//	b.Get("/like/{id:long}", GetLikeHandler)
//}

func IndexRoutes(m *mvc.Application)  {
	index := m.Party("/Index")
	index.Register(
		storage.Source,
		webApp.Sessions.Start,
	)
	index.Handle(new(indexController.IndexController))
}
