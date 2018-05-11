package reptilian

import (
	."github.com/alezh/novel/system/spider"
	"github.com/alezh/novel/system/http"
	"time"
	"github.com/alezh/novel/system/pipeline"
	"github.com/alezh/novel/system/http/request"
	"math/rand"
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
		pipeline.Pipeline              //结果收集与输出管道              //结果收集与输出管道
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
	self.Pipeline = pipeline.New(sp)
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
			//logs.Log.Debug(" *     Start: %v", req.GetUrl())
			self.Process(req)
		}()

		// 随机等待
		self.sleep()
	}

	// 等待处理中的任务完成
	self.Spider.Defer()
}

// core processer
func (self *reptilian) Process(req *request.Request) {
//	var (
//		downUrl = req.GetUrl()
//		sp      = self.Spider
//	)
//	defer func() {
//		if p := recover(); p != nil {
//			if sp.IsStopping() {
//				// println("Process$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
//				return
//			}
//			// 返回是否作为新的失败请求被添加至队列尾部
//			if sp.DoHistory(req, false) {
//				// 统计失败数
//				cache.PageFailCount()
//			}
//			// 提示错误
//			stack := make([]byte, 4<<10) //4KB
//			length := runtime.Stack(stack, true)
//			start := bytes.Index(stack, []byte("/src/runtime/panic.go"))
//			stack = stack[start:length]
//			start = bytes.Index(stack, []byte("\n")) + 1
//			stack = stack[start:]
//			if end := bytes.Index(stack, []byte("\ngoroutine ")); end != -1 {
//				stack = stack[:end]
//			}
//			stack = bytes.Replace(stack, []byte("\n"), []byte("\r\n"), -1)
//			logs.Log.Error(" *     Panic  [process][%s]: %s\r\n[TRACE]\r\n%s", downUrl, p, stack)
//		}
//	}()
//
//	var ctx = self.Downloader.Download(sp, req) // download page
//
//	if err := ctx.GetError(); err != nil {
//		// 返回是否作为新的失败请求被添加至队列尾部
//		if sp.DoHistory(req, false) {
//			// 统计失败数
//			cache.PageFailCount()
//		}
//		// 提示错误
//		logs.Log.Error(" *     Fail  [download][%v]: %v\n", downUrl, err)
//		return
//	}
//
//	// 过程处理，提炼数据
//	ctx.Parse(req.GetRuleName())
//
//	// 该条请求文件结果存入pipeline
//	for _, f := range ctx.PullFiles() {
//		if self.Pipeline.CollectFile(f) != nil {
//			break
//		}
//	}
//	// 该条请求文本结果存入pipeline
//	for _, item := range ctx.PullItems() {
//		if self.Pipeline.CollectData(item) != nil {
//			break
//		}
//	}
//
//	// 处理成功请求记录
//	sp.DoHistory(req, true)
//
//	// 统计成功页数
//	cache.PageSuccCount()
//
//	// 提示抓取成功
//	logs.Log.Informational(" *     Success: %v\n", downUrl)
//
//	// 释放ctx准备复用
//	spider.PutContext(ctx)
}


// 常用基础方法
func (self *reptilian) sleep() {
	sleeptime := self.pause[0] + rand.Int63n(self.pause[1])
	time.Sleep(time.Duration(sleeptime) * time.Millisecond)
}

// 从调度读取一个请求
func (self *reptilian) GetOne() *request.Request {
	return self.Spider.RequestPull()
}

//从调度使用一个资源空位
func (self *reptilian) UseOne() {
	self.Spider.RequestUse()
}

//从调度释放一个资源空位
func (self *reptilian) FreeOne() {
	self.Spider.RequestFree()
}

func (self *reptilian) SetId(id int) {
	self.id = id
}

func (self *reptilian) GetId() int {
	return self.id
}