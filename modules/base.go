package modules

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/alezh/novel/storage"
	"github.com/alezh/novel/system"
)

type BaseController struct {
	//mvc.BaseController
	Ctx iris.Context
	Source  *storage.DataSource
	Session *sessions.Session
	DbType string
}

func (c *BaseController) BeginRequest(ctx iris.Context) {
	if c.Source.Mysql != nil {
		system.SystemInfo.SetConfig("DbStype", "mysql")
	}else if c.Source.MongoDb != nil{
		system.SystemInfo.SetConfig("DbStype", "mgo")
	}
}

func (c *BaseController) EndRequest(ctx iris.Context) {
}
