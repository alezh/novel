package basics

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


//请求前处理
func (c *BaseController) BeginRequest(ctx iris.Context) {
}

//请求后处理
func (c *BaseController) EndRequest(ctx iris.Context) {
}
