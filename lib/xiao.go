package lib

import "github.com/alezh/novel/modules/spider"

func init()  {
	xiao.Register()
}

var xiao = &spider.Spider{
	Name:"xiao",
	Description:"book",
}
