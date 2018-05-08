package mission

import (
	"sync"
	"github.com/alezh/novel/config"
)

//任务调度与分配

type Mission struct {
	status       int          // 运行状态
	count        chan bool    // 总并发量计数
	useProxy     bool         // 标记是否使用代理IP
	//proxy        *proxy.Proxy // 全局代理IP
	matrices     []*Matrix    // Spider实例的请求矩阵列表
	sync.RWMutex              // 全局读写锁
}
// 定义全局调度
var mis = &Mission{
	status: config.RUN,
	count:  make(chan bool, config.Task.ThreadNum),
	//proxy:  proxy.New(),
}

//代理初始化
func Init() {
	//for mis.proxy == nil {
	//	time.Sleep(100 * time.Millisecond)
	//}
	mis.matrices = []*Matrix{}
	mis.count = make(chan bool, config.Task.ThreadNum)

	if config.Task.ProxyMinute > 0 {
	//	if mis.proxy.Count() > 0 {
	//		mis.useProxy = true
	//		mis.proxy.UpdateTicker(config.Task.ProxyMinute)
	//		//logs.Log.Informational(" *     使用代理IP，代理IP更换频率为 %v 分钟\n", cache.Task.ProxyMinute)
	//	} else {
	//		mis.useProxy = false
	//		//logs.Log.Informational(" *     在线代理IP列表为空，无法使用代理IP\n")
	//	}
	} else {
		mis.useProxy = false
		//logs.Log.Informational(" *     不使用代理IP\n")
	}
	mis.useProxy = false
	mis.status = config.RUN
}

// 注册资源队列
func AddMatrix(spiderName, spiderSubName string, maxPage int64) *Matrix {
	matrix := newMatrix(spiderName, spiderSubName, maxPage)
	mis.RLock()
	defer mis.RUnlock()
	mis.matrices = append(mis.matrices, matrix)
	return matrix
}

// 暂停\恢复所有爬行任务
func PauseRecover() {
	mis.Lock()
	defer mis.Unlock()
	switch mis.status {
	case config.PAUSE:
		mis.status = config.RUN
	case config.RUN:
		mis.status = config.PAUSE
	}
}

// 终止任务
func Stop() {
	mis.Lock()
	defer mis.Unlock()
	mis.status = config.STOP
	// 清空
	defer func() {
		recover()
	}()
	close(mis.count)
	mis.matrices = []*Matrix{}
}

// 每个spider实例分配到的平均资源量
func (self *Mission) avgRes() int32 {
	avg := int32(cap(mis.count) / len(mis.matrices))
	if avg == 0 {
		avg = 1
	}
	return avg
}

func (self *Mission) checkStatus(s int) bool {
	self.RLock()
	b := self.status == s
	self.RUnlock()
	return b
}