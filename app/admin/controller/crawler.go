package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/alezh/novel/system"
	"fmt"
)

// get:/Admin/crawler/list
func (c *AdminController)GetCrawlerList()mvc.Result{

	spiderMenu := func() (spmenu []map[string]string) {
		// 获取蜘蛛家族
		for _, sp := range system.SystemInfo.GetSpiderLib() {
			spmenu = append(spmenu, map[string]string{"name": sp.GetName(), "description": sp.GetDescription()})
		}
		return spmenu
	}()
	fmt.Println(spiderMenu)
	return mvc.View{
		Name:"admin/crawler/list.html",
		Data:spiderMenu,
	}
}