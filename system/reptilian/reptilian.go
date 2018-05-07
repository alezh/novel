package reptilian

import (
	."github.com/alezh/novel/system/spider"
	"github.com/alezh/novel/system/data/collector"
	"github.com/alezh/novel/system/http"
)

type (
	Reptilian interface {

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