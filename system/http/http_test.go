package http

import (
	"testing"
	"time"
	"github.com/alezh/novel/system/http/request"
	"github.com/alezh/novel/system/http/surfer"
	"io/ioutil"
	"github.com/alezh/novel/system/utils"
)

func TestHttpRequit(t *testing.T)  {
	t.Log(time.Now())
	req := &request.Request{
		DialTimeout: 5 * time.Second,
		Url:"http://www.biquke.com/bq/48/48586/",
		Method:"GET",
		TryTimes:5,
	}
	req.Prepare()
	resp, err :=surfer.Download(req)
	if err !=nil{
		t.Log(err)
	}
	b , errs := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if errs !=nil{
		t.Log(errs)
	}
	t.Log(utils.Bytes2String(b))
}