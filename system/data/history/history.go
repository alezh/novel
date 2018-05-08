package history

import (
	"sync"
	"github.com/alezh/novel/system/http/request"
)


type (
	Historier interface {
		//ReadSuccess(provider string, inherit bool) // 读取成功记录
		//UpsertSuccess(string) bool                 // 更新或加入成功记录
		HasSuccess(string) bool                    // 检查是否存在某条成功记录
		//DeleteSuccess(string)                      // 删除成功记录
		FlushSuccess(provider string)              // I/O输出成功记录，但不清缓存
		//
		//ReadFailure(provider string, inherit bool) // 取出失败记录
		//PullFailure() map[string]*request.Request  // 拉取失败记录并清空
		//UpsertFailure(*request.Request) bool       // 更新或加入失败记录
		//DeleteFailure(*request.Request)            // 删除失败记录
		FlushFailure(provider string)              // I/O输出失败记录，但不清缓存
	}
	History struct {
		*Success
		*Failure
		provider string
		sync.RWMutex
	}
)

func New(name string, subName string) Historier{
	return &History{
		Success: &Success{
			//tabName:  util.FileNameReplace(successTabName),
			//fileName: successFileName,
			new:      make(map[string]bool),
			old:      make(map[string]bool),
		},
		Failure: &Failure{
			//tabName:  util.FileNameReplace(failureTabName),
			//fileName: failureFileName,
			list:     make(map[string]*request.Request),
		},
	}
}


// I/O输出成功记录，但不清缓存
func (self *History) FlushSuccess(provider string) {
	self.RWMutex.Lock()
	self.provider = provider
	self.RWMutex.Unlock()
	sucLen, err := self.Success.flush(provider)
	if sucLen <= 0 {
		return
	}
	// logs.Log.Informational(" * ")
	if err != nil {
		//logs.Log.Error("%v", err)
	} else {
		//logs.Log.Informational(" *     [添加成功记录]: %v 条\n", sucLen)
	}
}
// I/O输出失败记录，但不清缓存
func (self *History) FlushFailure(provider string) {
	self.RWMutex.Lock()
	self.provider = provider
	self.RWMutex.Unlock()
	failLen, err := self.Failure.flush(provider)
	if failLen <= 0 {
		return
	}
	// logs.Log.Informational(" * ")
	if err != nil {
		//logs.Log.Error("%v", err)
	} else {
		//logs.Log.Informational(" *     [添加失败记录]: %v 条\n", failLen)
	}
}