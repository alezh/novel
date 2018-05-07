package http

import (
	"github.com/alezh/novel/system/http/surfer"
	"net/http/cookiejar"
	"net/http"
	"github.com/alezh/novel/system/spider"
	"github.com/alezh/novel/system/http/request"
)

type Surfer struct {
	surf    surfer.Surfer
	phantom surfer.Surfer
}
var (
	cookieJar, _     = cookiejar.New(nil)
	SurferDownloader = &Surfer{
		surf:    surfer.New(cookieJar),
	}
)

func (self *Surfer) Download(sp *spider.Spider, cReq *request.Request) *spider.Context {
	ctx := spider.GetContext(sp, cReq)

	var resp *http.Response
	var err error

	switch cReq.GetDownloaderID() {
	case request.SURF_ID:
		resp, err = self.surf.Download(cReq)

	case request.PHANTOM_ID:
		resp, err = self.phantom.Download(cReq)
	}

	if resp.StatusCode >= 400 {
		//err = errors.New("响应状态 " + resp.Status)
	}

	ctx.SetResponse(resp).SetError(err)

	return ctx
}