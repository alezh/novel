package reptilian

import (
	."github.com/alezh/novel/system/spider"
	"github.com/alezh/novel/system/pipeline/collector"
	"github.com/alezh/novel/system/http"
	"time"
)

type (
	Reptilian interface {
		Init(*Spider) Reptilian //初始化采集引擎
		Run()                        //运行任务
		Stop()                       //主动终止
		CanStop() bool               //能否终止
		GetId() int                  //获取引擎ID
	}
	reptilian struct {
		*Spider                 //执行的采集规则
		http.Downloader        //全局公用的下载器
		collector.Collector              //结果收集与输出管道
		id                    int      //引擎ID
		pause                 [2]int64 //[请求间隔的最短时长,请求间隔的增幅时长]
	}
)

func New(id int) Reptilian {
	return &reptilian{
		id:         id,
		Downloader: http.SurferDownloader,
	}
}

func (self *reptilian) Init(sp *Spider) Reptilian {
	self.Spider = sp.ReqmatrixInit()
	//随机暂停处理
	//self.Pipeline = pipeline.New(sp)
	//self.pause[0] = sp.Pausetime / 2
	//if self.pause[0] > 0 {
	//	self.pause[1] = self.pause[0] * 3
	//} else {
	//	self.pause[1] = 1
	//}
	self.pause[1] = 1
	return self
}

// 任务执行入口
func (self *reptilian) Run() {
	// 预先启动数据收集/输出管道
	self.Pipeline.Start()

	// 运行处理协程
	c := make(chan bool)
	go func() {
		self.run()
		close(c)
	}()

	// 启动任务
	self.Spider.Start()

	<-c // 等待处理协程退出

	// 停止数据收集/输出管道
	self.Pipeline.Stop()
}
// 主动终止
func (self *reptilian) Stop() {
	// 主动崩溃爬虫运行协程
	self.Spider.Stop()
	self.Pipeline.Stop()
}

func (self *reptilian) run() {
	for {
		// 队列中取出一条请求并处理
		req := self.GetOne()
		if req == nil {
			// 停止任务
			if self.Spider.CanStop() {
				break
			}
			time.Sleep(20 * time.Millisecond)
			continue
		}

		// 执行请求
		self.UseOne()
		go func() {
			defer func() {
				self.FreeOne()
			}()
			logs.Log.Debug(" *     Start: %v", req.GetUrl())
			self.Process(req)
		}()

		// 随机等待
		self.sleep()
	}

	// 等待处理中的任务完成
	self.Spider.Defer()
}