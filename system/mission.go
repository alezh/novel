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
	spider spider.Spider
	Request  *request.Request  //连接请求
	Response *http.Response    // 响应流
	text     []byte            // 下载内容Body的字节流格式
	dom      *goquery.Document // 下载内容Body为html时，可转换为Dom的对象
	err      error             // 错误标记
	Item     data.Data         //保存数据
	sync.Mutex
}
