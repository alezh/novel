package history

import "sync"

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
	//TODO::写入库操作
	return sLen,nil
}