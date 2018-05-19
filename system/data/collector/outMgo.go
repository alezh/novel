package collector

import (
	"fmt"
	"github.com/alezh/novel/system/utils"
)

func init() {
	DataOutput["mgo"] = saveMgo
}

func saveMgo(self *Collector) error {
	if self.MongoDb.Session.Ping() != nil{
		self.MongoDb.Session.Refresh()
		if self.MongoDb.Session.Ping() !=nil{
			fmt.Println("mgo 数据库连接错误")
		}
	}
	namespace := utils.FileNameReplace(self.Spider.GetName())
	collectlr :=  make(map[string][]interface{})
	for _, datacell := range self.dataDocker{
		subNamespace := utils.FileNameReplace(self.subNamespace(datacell))
		cName := joinNamespaces(namespace, subNamespace)
		saveDb := make(map[string]interface{})
		for k, v := range datacell["Data"].(map[string]interface{}) {
			saveDb[k] = v
		}
		collectlr[cName] = append(collectlr[cName],saveDb)
	}
	for table ,pdata := range collectlr{
		fmt.Println(len(pdata))
		self.MongoDb.InsetAll(table,pdata...)
	}
	return nil
}