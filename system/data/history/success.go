package history

import (
	"sync"
	"fmt"
)

//成功记录保存
type Success struct {
	tabName     string
	fileName    string
	new         map[string]bool // [Request.Unique()]true
	old         map[string]bool // [Request.Unique()]true
	inheritable bool
	sync.RWMutex
}

func (self *Success) HasSuccess(reqUnique string) bool {
	self.RWMutex.Lock()
	has := self.old[reqUnique] || self.new[reqUnique]
	self.RWMutex.Unlock()
	return has
}

func (self *Success) flush(provider string) (sLen int, err error){
	self.RWMutex.Lock()
	defer self.RWMutex.Unlock()

	sLen = len(self.new)
	if sLen == 0 {
		return
	}
	switch provider{
	case "mgo":
		fmt.Println("成功记录写入库操作")
		//storage.Source.MongoDb.Database.C("novel").Count()
	case "mysql":
		//同步表
		//storage.Source.Mysql.Sync2(&FailureModel{})
	}
	//TODO::写入库操作
	return sLen,nil
}