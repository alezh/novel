package lib

import (
	. "github.com/alezh/novel/system/spider"
	"github.com/alezh/novel/system/http/request"
	"github.com/PuerkitoBio/goquery"
)

func init()  {
	test.Register()
}

var test = &Spider{
	Name:        "抓取测试",
	Description: "抓取测试",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{
				Url:  "https://www.biquge5200.cc/79_79875/",
				Rule: "列表",
			})
		},
		Trunk: map[string]*Rule{
			"列表": {
				ItemFields: []string{
					"Title",
					"Url",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					id := "#list dd a"
					query.Find(id).Each(func(j int, s *goquery.Selection){
						url, _ := s.Attr("href")
						text := s.Text()
						// 结果存入Response中转
						ctx.Output(map[int]interface{}{
							0: text,
							1: url,
						})
					})

				},
			},
		},
	},
}