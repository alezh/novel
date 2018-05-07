package http

import (
	"github.com/alezh/novel/system/spider"
	"github.com/alezh/novel/system/http/request"
)

type Downloader interface {
	Download(*spider.Spider, *request.Request)*spider.Context
} 