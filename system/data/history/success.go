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
