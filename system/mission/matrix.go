package mission

import (
	"sync"
	"github.com/alezh/novel/system/http/request"
	"github.com/alezh/novel/config"
)

type Matrix struct {
	maxPage         int64                       // 最大采集页数，以负数形式表示
	resCount        int32                       // 资源使用情况计数
	spiderName      string                      // 所属Spider
	reqs            map[int][]*request.Request  // [优先级]队列，优先级默认为0
	priorities      []int                       // 优先级顺序，从低到高
	//history         history.Historier           // 历史记录
	tempHistory     map[string]bool             // 临时记录 [reqUnique(url+method)]true
	failures        map[string]*request.Request // 历史及本次失败请求
	tempHistoryLock sync.RWMutex
	failureLock     sync.Mutex
	sync.Mutex
}

func newMatrix(spiderName, spiderSubName string, maxPage int64) *Matrix {
	matrix := &Matrix{
		spiderName:  spiderName,
		maxPage:     maxPage,
		reqs:        make(map[int][]*request.Request),
		priorities:  []int{},
		//history:     history.New(spiderName, spiderSubName),
		tempHistory: make(map[string]bool),
		failures:    make(map[string]*request.Request),
	}
	if config.Task.Mode != config.SERVER {
		//matrix.history.ReadSuccess(config.SysConfig..OutType, config.SysConfig..SuccessInherit)
		//matrix.history.ReadFailure(config.SysConfig..OutType, config.SysConfig..FailureInherit)
		//matrix.setFailures(matrix.history.PullFailure())
	}
	return matrix
}


