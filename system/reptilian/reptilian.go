package reptilian

import (
	."github.com/alezh/novel/system/spider"
	"github.com/alezh/novel/system/data/collector"
	"github.com/alezh/novel/system/http"
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
