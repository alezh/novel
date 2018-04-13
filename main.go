package main

import (
	"github.com/alezh/novel/app"
	"github.com/kataras/iris"
)

func main()  {
	app.NewApp().Run(iris.Addr(":8080"),app.Configure)
}

