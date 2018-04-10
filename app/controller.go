package app

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/alezh/novel/storage"
)

type BaseController struct {
	Ctx iris.Context
	Source  *storage.DataSource
	Session *sessions.Session
}