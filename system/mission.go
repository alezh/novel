package system

import (
	"github.com/alezh/novel/modules/spider"
	"github.com/alezh/novel/system/http/request"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"sync"
	"github.com/alezh/novel/system/data"
)

//任务调度与分配

type Mission struct {
	spider   *spider.Spider
	Request  *request.Request  //连接请求
	Response *http.Response    // 响应流
	text     []byte            // 下载内容Body的字节流格式
	dom      *goquery.Document // 下载内容Body为html时，可转换为Dom的对象
	err      error             // 错误标记
	item     []data.Data         // 保存数据
	sync.Mutex
}

var missionPool = sync.Pool{
	New: func() interface{} {
		return &Mission{
			item:[]data.Data{},
		}
	},
}

//********************GET/PUT*************************************************/

func Get(sp *spider.Spider, req *request.Request) *Mission {
	mission :=missionPool.Get().(*Mission)
	mission.spider = sp
	mission.Request = req
	return mission
}

func Put(miss *Mission)  {
	if miss.Response != nil {
		miss.Response.Body.Close() // too many open files bug remove
		miss.Response = nil
	}
	miss.item = miss.item[:0]
	miss.spider = nil
	miss.Request = nil
	miss.text = nil
	miss.dom = nil
	miss.err = nil
	missionPool.Put(miss)
}

func (miss *Mission) SetResponse(resp *http.Response) *Mission {
	miss.Response = resp
	return miss
}

func (miss *Mission) AddQueue(req *request.Request) *Mission{
	miss.spider.tryPanic()

	return miss
}
