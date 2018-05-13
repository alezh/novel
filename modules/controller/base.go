package controller

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/alezh/novel/storage"
)

type BaseController struct {
	Ctx iris.Context
	Source  *storage.DataSource
	Session *sessions.Session
	DbType string
}

//var dbType = system.SystemInfo.GetConfig("DbStype").(string)

//请求前处理
func (c *BaseController) BeginRequest(ctx iris.Context) {
	//if dbType == ""{
	//	if c.Source.Mysql != nil {
	//		system.SystemInfo.SetConfig("DbStype", "mysql")
	//		dbType = "mysql"
	//	}else if c.Source.MongoDb != nil{
	//		system.SystemInfo.SetConfig("DbStype", "mgo")
	//		dbType = "mysql"
	//	}
	//}
}

//请求后处理
func (c *BaseController) EndRequest(ctx iris.Context) {
}
