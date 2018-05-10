package collector

import (
	"fmt"
	"github.com/alezh/novel/system/utils"
)

/************************ MongoDB 输出 ***************************/

func init() {
	DataOutput["mgo"] = func(self *Collector) error {
		//连接数据库
		if self.MongoDb == nil {
			return fmt.Println("MongoBD数据库链接失败")
		}

		return mgo.Call(func(src pool.Src) error {
			var (
				db          = src.(*mgo.MgoSrc).DB(config.DB_NAME)
				namespace   = utils.FileNameReplace(self.namespace())
				collections = make(map[string]*mgov2.Collection)
				dataMap     = make(map[string][]interface{})
				err         error
			)

			for _, datacell := range self.dataDocker {
				subNamespace := utils.FileNameReplace(self.subNamespace(datacell))
				cName := joinNamespaces(namespace, subNamespace)

				if _, ok := collections[subNamespace]; !ok {
					collections[subNamespace] = db.C(cName)
				}
				for k, v := range datacell["Data"].(map[string]interface{}) {
					datacell[k] = v
				}
				delete(datacell, "Data")
				delete(datacell, "RuleName")
				if !self.Spider.OutDefaultField() {
					delete(datacell, "Url")
					delete(datacell, "ParentUrl")
					delete(datacell, "DownloadTime")
				}
				dataMap[subNamespace] = append(dataMap[subNamespace], datacell)
			}

			for collection, docs := range dataMap {
				c := collections[collection]
				count := len(docs)
				loop := count / mgo.MaxLen
				for i := 0; i < loop; i++ {
					err = c.Insert(docs[i*mgo.MaxLen : (i+1)*mgo.MaxLen]...)
					if err != nil {
						logs.Log.Error("%v", err)
					}
				}
				if count%mgo.MaxLen == 0 {
					continue
				}
				err = c.Insert(docs[loop*mgo.MaxLen:]...)
				if err != nil {
					logs.Log.Error("%v", err)
				}
			}

			return nil
		})
	}
}
