package lib

import "github.com/alezh/novel/system/spider"

func init()  {
	xiao.Register()
}

var xiao = &spider.Spider{
	Name:"xiao",
	Description:"book",
}
