package main

import (
	"github.com/alezh/novel/app"
	_ "github.com/alezh/novel/lib"
	"github.com/kataras/iris"
)

func main()  {
	app.NewApp().Run(iris.Addr(":8080"),app.Configure)
}

