package controller

import (
	"github.com/kataras/iris/mvc"
	"github.com/alezh/novel/system"
	"github.com/kataras/iris"
	"github.com/alezh/novel/system/utils"
	"github.com/alezh/novel/config"
	"github.com/alezh/novel/system/spider"
	"encoding/json"
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
	return mvc.View{
		Name:"admin/crawler/list.html",
		Data:iris.Map{"data":spiderMenu},
	}
}

// 添加任务 post:/Admin/add/queue
func (c *AdminController)PostAddQueue(form formValue)  {
	var spNames map[string]interface{}
	data := form("spider")
	if data != ""{
		json.Unmarshal([]byte(data), &spNames)
	}
	var mode = utils.Atoi(form("mode"))
	var port = utils.Atoi(form("port"))
	var master = utils.Atoa(form("ip")) //服务器(主节点)地址，不含端口
	currMode := system.SystemInfo.GetConfig("mode").(int)
	if currMode == config.UNSET{
		system.SystemInfo.Init(mode, port, master)
	}else{
		system.SystemInfo.ReInit(mode, port, master)
	}

	spiders := []*spider.Spider{}
	sp := system.SystemInfo.GetSpiderByName("抓取测试")
	spiders = append(spiders, sp.Copy())
	//for _, sp := range system.SystemInfo.GetSpiderLib() {
	//	for _, spName := range spNames {
	//		if utils.Atoa(spName) == sp.GetName() {
	//			spiders = append(spiders, sp.Copy())
	//		}
	//	}
	//}
	system.SystemInfo.SpiderPrepare(spiders)
	jsons := iris.Map{"len":system.SystemInfo.GetSpiderQueue().Len()}
	c.Ctx.JSON(jsons)
}




// 服务器开启 POST :/Admin/engine/start
func (c *AdminController)PostEngineStart(form formValue)  {
	system.SystemInfo.SetConfig("Mode", config.OFFLINE)

	go func() {
		system.SystemInfo.Start()
		//if system.SystemInfo.GetConfig("mode").(int) == config.OFFLINE {
			//Sc.Write(sessID, map[string]interface{}{"operate": "stop"})
		//}
	}()
	c.Ctx.WriteString("go")
}

func (c *AdminController)PostEngineStop(form formValue)  {

}