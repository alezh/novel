package reptilian

import (
	"sync"
	"github.com/alezh/novel/config"
	"time"
)

type (
	ReptilianPool interface {
		Reset(int) int
		Use() Reptilian
		Free(Reptilian)
		Stop()
	}
	rp struct {
		capacity int
		count    int
		usable   chan Reptilian
		all      []Reptilian
		status   int
		sync.RWMutex
	}
)

func NewReptilianPool() ReptilianPool {
	return &rp{
		status: config.RUN,
		all:    make([]Reptilian, 0, config.CRAWLS_CAP), // 蜘蛛池最大容量
	}
}

// 根据要执行的蜘蛛数量设置CrawlerPool
// 在二次使用Pool实例时，可根据容量高效转换
func (self *rp) Reset(spiderNum int) int {
	self.Lock()
	defer self.Unlock()
	var wantNum int
	if spiderNum < config.CRAWLS_CAP {
		wantNum = spiderNum
	} else {
		wantNum = config.CRAWLS_CAP
	}
	if wantNum <= 0 {
		wantNum = 1
	}
	self.capacity = wantNum
	self.count = 0
	self.usable = make(chan Reptilian, wantNum)
	for _, crawler := range self.all {
		if self.count < self.capacity {
			self.usable <- crawler
			self.count++
		}
	}
	self.status = config.RUN
	return wantNum
}

// 并发安全地使用资源
func (self *rp) Use() Reptilian {
	var crawler Reptilian
	for {
		self.Lock()
		if self.status == config.STOP {
			self.Unlock()
			return nil
		}
		select {
		case crawler = <-self.usable:
			self.Unlock()
			return crawler
		default:
			if self.count < self.capacity {
				crawler = New(self.count)
				self.all = append(self.all, crawler)
				self.count++
				self.Unlock()
				return crawler
			}
		}
		self.Unlock()
		time.Sleep(time.Second)
	}
	return nil
}

func (self *rp) Free(crawler Reptilian) {
	self.RLock()
	defer self.RUnlock()
	if self.status == config.STOP || !crawler.CanStop() {
		return
	}
	self.usable <- crawler
}

// 主动终止所有爬行任务
func (self *rp) Stop() {
	self.Lock()
	// println("CrawlerPool^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	if self.status == config.STOP {
		self.Unlock()
		return
	}
	self.status = config.STOP
	close(self.usable)
	self.usable = nil
	self.Unlock()

	for _, crawler := range self.all {
		crawler.Stop()
	}
	// println("CrawlerPool$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
}