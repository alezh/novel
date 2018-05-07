package collector

import (
	"sync"
	."github.com/alezh/novel/system/spider"
	."github.com/alezh/novel/system/data"
)

type Collector struct {
	*Spider                    //绑定的采集规则
	DataChan       chan Data //文本数据收集通道
	dataDocker     []Data    //分批输出结果缓存
	outType        string             //输出方式
	// size     [2]uint64 //数据总输出流量统计[文本，文件]，文本暂时未统计
	dataBatch   uint64 //当前文本输出批次
	fileBatch   uint64 //当前文件输出批次
	wait        sync.WaitGroup
	sum         [4]uint64 //收集的数据总数[上次输出后文本总数，本次输出后文本总数，上次输出后文件总数，本次输出后文件总数]，非并发安全
	dataSumLock sync.RWMutex
	fileSumLock sync.RWMutex
}

func NewCollector(sp *Spider) *Collector {
	var self = &Collector{}
	self.Spider = sp
	//self.outType = cache.Task.OutType
	//if cache.Task.DockerCap < 1 {
	//	cache.Task.DockerCap = 1
	//}
	//self.DataChan = make(chan data.DataCell, cache.Task.DockerCap)
	//self.FileChan = make(chan data.FileCell, cache.Task.DockerCap)
	//self.dataDocker = make([]data.DataCell, 0, cache.Task.DockerCap)
	self.sum = [4]uint64{}
	// self.size = [2]uint64{}
	self.dataBatch = 0
	self.fileBatch = 0
	return self
}