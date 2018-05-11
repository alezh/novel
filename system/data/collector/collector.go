package collector

import (
	"sync"
	."github.com/alezh/novel/system/spider"
	."github.com/alezh/novel/system/data"
	"fmt"
	"github.com/alezh/novel/config"
	"time"
)

type Collector struct {
	*Spider                    //绑定的采集规则
	//*storage.DataSource
	DataChan       chan Data //文本数据收集通道
	dataDocker     []Data    //分批输出结果缓存
	outType        string             //输出方式
	size     [2]uint64 //数据总输出流量统计[文本，文件]，文本暂时未统计
	dataBatch   uint64 //当前文本输出批次
	//fileBatch   uint64 //当前文件输出批次
	wait        sync.WaitGroup
	sum         [4]uint64 //收集的数据总数[上次输出后文本总数，本次输出后文本总数，上次输出后文件总数，本次输出后文件总数]，非并发安全
	dataSumLock sync.RWMutex
	//fileSumLock sync.RWMutex
}

func NewCollector(sp *Spider) *Collector {
	var self = &Collector{}
	self.Spider = sp
	//self.DataSource = storage.Source
	self.outType = config.Task.OutType
	if config.Task.DockerCap < 1 {
		config.Task.DockerCap = 1
	}
	self.DataChan = make(chan Data, config.Task.DockerCap)
	//self.FileChan = make(chan data.FileCell, cache.Task.DockerCap)
	self.dataDocker = make([]Data, 0, config.Task.DockerCap)
	self.sum = [4]uint64{}
	self.size = [2]uint64{}
	self.dataBatch = 0
	//self.fileBatch = 0
	return self
}

func (self *Collector) CollectData(dataCell Data) error {
	var err error
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("输出协程已终止")
		}
	}()
	self.DataChan <- dataCell
	return err
}

// 停止
func (self *Collector) Stop() {
	go func() {
		defer func() {
			recover()
		}()
		close(self.DataChan)
	}()
}
// 启动数据收集/输出管道
func (self *Collector) Start() {
	// 启动输出协程
	go func() {
		dataStop := make(chan bool)
		fileStop := make(chan bool)

		go func() {
			defer func() {
				recover()
				// println("DataChanStop$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
			}()
			for data := range self.DataChan {
				// 缓存分批数据
				self.dataDocker = append(self.dataDocker, data)

				// 未达到设定的分批量时继续收集数据
				if len(self.dataDocker) < config.Task.DockerCap {
					continue
				}

				// 执行输出
				self.dataBatch++
				self.outputData()
			}
			// 将剩余收集到但未输出的数据输出
			self.dataBatch++
			self.outputData()
			close(dataStop)
		}()

		go func() {
			defer func() {
				recover()
				// println("FileChanStop$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
			}()
			//// 只有当收到退出通知并且通道内无数据时，才退出循环
			//for file := range self.FileChan {
			//	atomic.AddUint64(&self.fileBatch, 1)
			//	self.wait.Add(1)
			//	go self.outputFile(file)
			//}
			close(fileStop)
		}()

		<-dataStop
		<-fileStop
		// println("OutputWaitStopping++++++++++++++++++++++++++++++++")

		// 等待所有输出完成
		//self.wait.Wait()
		// println("OutputStopped$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")

		// 返回报告
		self.Report()
	}()
}
func (self *Collector) resetDataDocker() {
	for _, cell := range self.dataDocker {
		PutDataCell(cell)
	}
	self.dataDocker = self.dataDocker[:0]
}

// 获取文本数据总量
func (self *Collector) dataSum() uint64 {
	self.dataSumLock.RLock()
	defer self.dataSumLock.RUnlock()
	return self.sum[1]
}

// 更新文本数据总量
func (self *Collector) addDataSum(add uint64) {
	self.dataSumLock.Lock()
	defer self.dataSumLock.Unlock()
	self.sum[0] = self.sum[1]
	self.sum[1] += add
}
// 返回报告
func (self *Collector) Report() {
	config.ReportChan <- &config.Report{
		SpiderName: self.Spider.GetName(),
		Keyin:      self.GetKeyin(),
		DataNum:    self.dataSum(),
		// DataSize:   self.dataSize(),
		// FileSize: self.fileSize(),
		Time: time.Since(config.StartTime),
	}
}