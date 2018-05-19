package history

import (
	"sync"
	"github.com/alezh/novel/system/http/request"
	"time"
	"fmt"
)
//失败纪录保保存
type (
	Failure struct {
		tabName     string
		fileName    string
		list        map[string]*request.Request //key:url
		inheritable bool
		sync.RWMutex
	}
	FailureModel struct {
		Id        int64       `xorm:"not null pk autoincr BIGINT(20)"`
		Key       string      `xorm:"VARCHAR(255) NOT NULL PRIMARY KEY"`
		Failure   string      `xorm:"MEDIUMTEXT"`
		CreatedAt time.Time   `xorm:"created"`
	}

)

//去除错误记录
func (self *Failure) PullFailure() map[string]*request.Request {
	list := self.list
	self.list = make(map[string]*request.Request)
	return list
}

// 更新或加入失败记录，
// 对比是否已存在，不存在就记录，
// 返回值表示是否有插入操作。
func (self *Failure) UpsertFailure(req *request.Request) bool {
	self.RWMutex.Lock()
	defer self.RWMutex.Unlock()
	if self.list[req.Unique()] != nil {
		return false
	}
	self.list[req.Unique()] = req
	return true
}

// 删除失败记录
func (self *Failure) DeleteFailure(req *request.Request) {
	self.RWMutex.Lock()
	delete(self.list, req.Unique())
	self.RWMutex.Unlock()
}

func (self *Failure) flush(provider string) (fLen int, err error){
	self.RWMutex.Lock()
	defer self.RWMutex.Unlock()
	fLen = len(self.list)
	//TODO::写入库操作
	switch provider{
	case "mgo":
		fmt.Println("错误记录写入库操作",self.list)
		//storage.Source.MongoDb.Database.C("novel").Count()
	case "mysql":
		//同步表
		//storage.Source.Mysql.Sync2(&FailureModel{})
	}
	return fLen ,nil
}