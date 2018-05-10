// 数据收集
package data

import (
	."github.com/alezh/novel/system/spider"
	"github.com/alezh/novel/system/data/collector"
)

// 数据收集/输出管道
type Pipeline interface {
	Start()                          //启动
	Stop()                           //停止
	CollectData(Data) error          //收集数据单元
}

func New(sp *Spider) Pipeline {
	return collector.NewCollector(sp)
}
