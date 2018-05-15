package system

import (
	"github.com/henrylee2cn/teleport"
	"time"
	"sync"
	"github.com/alezh/novel/system/spider"
	"github.com/alezh/novel/config"
	"io"
	"github.com/alezh/novel/system/mission"
	"reflect"
	"strings"
	"math"
	"github.com/alezh/novel/system/reptilian"
	"strconv"
	"github.com/alezh/novel/system/data/collector"
	"fmt"
	"github.com/alezh/novel/system/pipeline"
	"runtime"
)

func init()  {
	// 开启最大核心数运行
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 开启手动GC
	ManualGC()
}

type (
	Enging interface {
		Init(int,int, string, ...io.Writer) Enging
		ReInit(int,int, string, ...io.Writer) Enging
		SpiderPrepare([]*spider.Spider) Enging
		GetOutputLib() []string
		GetSpiderLib() []*spider.Spider
		GetSpiderByName(string) *spider.Spider
		GetTaskJar() *mission.TaskJar
		GetSpiderQueue() reptilian.SpiderQueue
		Start()
		SetConfig(string,interface{})Enging
		GetConfig(...string)interface{}
	}
	System struct {
		*config.SysConfig                   // 全局配置
		*spider.Species                     // 全部蜘蛛种类
		*mission.TaskJar                    // 服务器与客户端间传递任务的存储库
		reptilian.SpiderQueue               // 当前任务的蜘蛛队列
		teleport.Teleport                   // socket长连接双工通信接口，json数据传输
		reptilian.ReptilianPool             // 爬行回收池
		sum                   [2]uint64     // 执行计数
		takeTime              time.Duration // 执行计时
		status                int           // 运行状态
		finish                chan bool
		finishOnce            sync.Once
		canSocketLog          bool
		sync.RWMutex
	}
)

var SystemInfo  = SysInterface()

func SysInterface() Enging {
	return initSystem()
}

func initSystem() *System {
	sys := &System{
		SysConfig:   config.Task,
		Species:     spider.SpeciesCollection,
		status:      config.STOPPED,
		SpiderQueue: reptilian.NewQueue(),
		Teleport:    teleport.New(),
		TaskJar:     mission.NewTaskJar(),
		ReptilianPool:reptilian.NewReptilianPool(),
	}
	return sys
}
//
//func WebStart()  {
//	app.NewApp().Run(iris.Addr(":8080"),app.Configure)
//}

func (sys *System)Init(mode int, port int, master string, w ...io.Writer) Enging {
	sys.canSocketLog = false
	if len(w) > 0 {
		//sys.SetLog(w[0])
	}
	//sys.LogGoOn()
	sys.SysConfig.Mode, sys.SysConfig.Port, sys.SysConfig.Master = mode, port, master
	sys.Teleport = teleport.New()
	sys.TaskJar = mission.NewTaskJar()
	sys.SpiderQueue = reptilian.NewQueue()

	switch sys.SysConfig.Mode{
	case config.SERVER:
		if sys.checkPort() {
			//logs.Log.Informational("                                                                                               ！！当前运行模式为：[ 服务器 ] 模式！！")
			sys.Teleport.SetAPI(mission.MasterApi(sys)).Server(":" + strconv.Itoa(sys.SysConfig.Port))
		}
	case config.CLIENT:
	case config.OFFLINE:
		return sys
	default:
		return sys
	}
	return sys
}

func (sys *System)ReInit(mode int, port int, master string, w ...io.Writer)  Enging{
	if !sys.IsStopped() {
		sys.Stop()
	}
	//sys.LogRest()
	if sys.Teleport != nil {
		sys.Teleport.Close()
	}
	// 等待结束
	if mode == config.UNSET {
		sys = initSystem()
		sys.SysConfig.Mode = config.UNSET
		return sys
	}
	// 重新开启
	sys = initSystem().Init(mode, port, master, w...).(*System)
	return sys
}

// SpiderPrepare()必须在设置全局运行参数之后，Run()的前一刻执行
// original为spider包中未有过赋值操作的原始蜘蛛种类
// 已被显式赋值过的spider将不再重新分配Keyin
// client模式下不调用该方法
func (self *System) SpiderPrepare(original []*spider.Spider) Enging {
	self.SpiderQueue.Reset()
	// 遍历任务
	for _, sp := range original {
		spcopy := sp.Copy()
		//spcopy.SetPausetime(self.SysConfig.Pausetime)
		if spcopy.GetLimit() == spider.LIMIT {
			spcopy.SetLimit(self.SysConfig.Limit)
		} else {
			spcopy.SetLimit(-1 * self.SysConfig.Limit)
		}
		self.SpiderQueue.Add(spcopy)
	}
	// 遍历自定义配置
	self.SpiderQueue.AddKeyins(self.SysConfig.Keyins)
	return self
}

// 获取全部输出方式
func (self *System) GetOutputLib() []string {
	return collector.DataOutputLib
}

// 获取全部蜘蛛种类
func (self *System) GetSpiderLib() []*spider.Spider {
	return self.Species.Get()
}

// 通过名字获取某蜘蛛
func (self *System) GetSpiderByName(name string) *spider.Spider {
	return self.Species.GetByName(name)
}

// 返回当前运行模式
func (self *System) GetMode() int {
	return self.SysConfig.Mode
}

// 返回任务库
func (self *System) GetTaskJar() *mission.TaskJar {
	return self.TaskJar
}
// 服务器客户端模式下返回节点数
func (self *System) CountNodes() int {
	return self.Teleport.CountNodes()
}

// 获取蜘蛛队列接口实例
func (self *System) GetSpiderQueue() reptilian.SpiderQueue {
	return self.SpiderQueue
}

// 运行任务
func (sys *System)Start()  {
	// 确保开启报告
	sys.LogGoOn()
	if sys.SysConfig.Mode != config.CLIENT && sys.SpiderQueue.Len() == 0 {
		//logs.Log.Warning(" *     —— 亲，任务列表不能为空哦~")
		//sys.LogRest()
		return
	}
	sys.finish = make(chan bool)
	sys.finishOnce = sync.Once{}
	// 重置计数
	sys.sum[0], sys.sum[1] = 0, 0
	// 重置计时
	sys.takeTime = 0
	// 设置状态
	sys.setStatus(config.RUN)
	defer sys.setStatus(config.STOPPED)
	// 任务执行
	switch sys.SysConfig.Mode {
	case config.OFFLINE:
		sys.offline()
	case config.SERVER:
		sys.server()
	case config.CLIENT:
		sys.client()
	default:
		return
	}
	<-sys.finish
}

func (sys *System)Stop()  {

}

func (sys *System)GetConfig(k ...string) interface{}  {
	defer func() {
		if err := recover(); err != nil {
			//logs.Log.Error("%v", err)
		}
	}()
	if len(k) == 0 {
		return sys.SysConfig
	}
	key := strings.Title(k[0])
	acv := reflect.ValueOf(sys.SysConfig).Elem()
	return acv.FieldByName(key).Interface()
}

// 设置全局参数
func (sys *System)SetConfig(k string, v interface{}) Enging {
	defer func() {
		if err := recover(); err != nil {
			//logs.Log.Error("%v", err)
		}
	}()
	if k == "Limit" && v.(int64) <= 0 {
		v = int64(math.MaxInt64) //spider.LIMIT
	} else if k == "DockerCap" && v.(int) < 1 {
		v = int(1)
	}
	//反射
	acv := reflect.ValueOf(sys.SysConfig).Elem()
	key := strings.Title(k)
	if acv.FieldByName(key).CanSet() {
		acv.FieldByName(key).Set(reflect.ValueOf(v))
	}
	return sys
}

func (sys *System)LogGoOn()  {
	
}

// 检查任务是否正在运行
func (self *System) IsRunning() bool {
	return self.status == config.RUN
}

// 检查任务是否处于暂停状态
func (self *System) IsPause() bool {
	return self.status == config.PAUSE
}

// 检查任务是否已经终止
func (self *System) IsStopped() bool {
	return self.status == config.STOPPED
}
// 返回当前运行状态
func (self *System) Status() int {
	self.RWMutex.RLock()
	defer self.RWMutex.RUnlock()
	return self.status
}


// 返回当前运行状态
func (self *System) setStatus(status int) {
	self.RWMutex.Lock()
	defer self.RWMutex.Unlock()
	self.status = status
}

// ******************************************** 运行方法 ************************************************* \\
// 离线模式运行
func (self *System) offline() {
	self.exec()
}

// 服务器模式运行，必须在SpiderPrepare()执行之后调用才可以成功添加任务
// 生成的任务与自身当前全局配置相同
func (self *System) server() {
	// 标记结束
	defer func() {
		self.finishOnce.Do(func() { close(self.finish) })
	}()

	// 便利添加任务到库
	tasksNum, spidersNum := self.addNewTask()

	if tasksNum == 0 {
		return
	}
	fmt.Printf("—— 本次成功添加 %v 条任务，共包含 %v 条采集规则 ——", tasksNum, spidersNum)
	// 打印报告
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(` *********************************************************************************************************************************** `)
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(" *                               —— 本次成功添加 %v 条任务，共包含 %v 条采集规则 ——", tasksNum, spidersNum)
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(` *********************************************************************************************************************************** `)
}

// 服务器模式下，生成task并添加至库
func (self *System) addNewTask() (tasksNum, spidersNum int) {
	length := self.SpiderQueue.Len()
	t := mission.Task{}
	// 从配置读取字段
	self.setTask(&t)
	for i, sp := range self.SpiderQueue.GetAll() {
		t.Spiders = append(t.Spiders, map[string]string{"name": sp.GetName(), "keyin": sp.GetKeyin()})
		spidersNum++
		// 每十个蜘蛛存为一个任务
		if i > 0 && i%10 == 0 && length > 10 {
			// 存入
			one := t
			self.TaskJar.Push(&one)
			// logs.Log.App(" *     [新增任务]   详情： %#v", *t)

			tasksNum++

			// 清空spider
			t.Spiders = []map[string]string{}
		}
	}
	if len(t.Spiders) != 0 {
		// 存入
		one := t
		self.TaskJar.Push(&one)
		tasksNum++
	}
	return
}

// 客户端模式运行
func (self *System) client() {
	// 标记结束
	defer func() {
		self.finishOnce.Do(func() { close(self.finish) })
	}()

	for {
		// 从任务库获取一个任务
		t := self.downTask()

		if self.Status() == config.STOP || self.Status() == config.STOPPED {
			return
		}

		// 准备运行
		self.taskToRun(t)

		// 重置计数
		self.sum[0], self.sum[1] = 0, 0
		// 重置计时
		self.takeTime = 0

		// 执行任务
		self.exec()
	}
}

// 客户端模式下获取任务
func (self *System) downTask() *mission.Task {
ReStartLabel:
	if self.Status() == config.STOP || self.Status() == config.STOPPED {
		return nil
	}
	if self.CountNodes() == 0 && self.TaskJar.Len() == 0 {
		time.Sleep(time.Second)
		goto ReStartLabel
	}

	if self.TaskJar.Len() == 0 {
		self.Request(nil, "task", "")
		for self.TaskJar.Len() == 0 {
			if self.CountNodes() == 0 {
				goto ReStartLabel
			}
			time.Sleep(time.Second)
		}
	}
	return self.TaskJar.Pull()
}

// client模式下从task准备运行条件
func (self *System) taskToRun(t *mission.Task) {
	// 清空历史任务
	self.SpiderQueue.Reset()

	// 更改全局配置
	self.setAppConf(t)

	// 初始化蜘蛛队列
	for _, n := range t.Spiders {
		sp := self.GetSpiderByName(n["name"])
		if sp == nil {
			continue
		}
		spcopy := sp.Copy()
		//spcopy.SetPausetime(t.Pausetime)
		if spcopy.GetLimit() > 0 {
			spcopy.SetLimit(t.Limit)
		} else {
			spcopy.SetLimit(-1 * t.Limit)
		}
		if v, ok := n["keyin"]; ok {
			spcopy.SetKeyin(v)
		}
		self.SpiderQueue.Add(spcopy)
	}
}

// 开始执行任务
func (self *System) exec() {
	count := self.SpiderQueue.Len()
	config.ResetPageCount()
	// 刷新输出方式的状态
	pipeline.RefreshOutput()
	// 初始化资源队列
	mission.Init()

	// 设置爬虫队列
	crawlerCap := self.ReptilianPool.Reset(count)
	fmt.Println(" *     采集引擎池容量为 :", crawlerCap)
	//logs.Log.Informational(" *     执行任务总数(任务数[*自定义配置数])为 %v 个\n", count)
	//logs.Log.Informational(" *     采集引擎池容量为 %v\n", crawlerCap)
	//logs.Log.Informational(" *     并发协程最多 %v 个\n", self.AppConf.ThreadNum)
	//logs.Log.Informational(" *     默认随机停顿 %v~%v 毫秒\n", self.AppConf.Pausetime/2, self.AppConf.Pausetime*2)
	//logs.Log.App(" *                                                                                                 —— 开始抓取，请耐心等候 ——")
	//logs.Log.Informational(` *********************************************************************************************************************************** `)

	// 开始计时
	config.StartTime = time.Now()

	// 根据模式选择合理的并发
	if self.SysConfig.Mode == config.OFFLINE {
		// 可控制执行状态
		go self.goRun(count)
	} else {
		// 保证接收服务端任务的同步
		self.goRun(count)
	}
}

// 任务执行
func (self *System) goRun(count int) {
	// 执行任务
	var i int
	for i = 0; i < count && self.Status() != config.STOP; i++ {
	pause:
		if self.IsPause() {
			time.Sleep(time.Second)
			goto pause
		}
		// 从爬行队列取出空闲蜘蛛，并发执行
		c := self.ReptilianPool.Use()
		if c != nil {
			go func(i int, c reptilian.Reptilian) {
				// 执行并返回结果消息
				c.Init(self.SpiderQueue.GetByIndex(i)).Run()
				// 任务结束后回收该蜘蛛
				self.RWMutex.RLock()
				if self.status != config.STOP {
					self.ReptilianPool.Free(c)
				}
				self.RWMutex.RUnlock()
			}(i, c)
		}
	}
	// 监控结束任务
	for ii := 0; ii < i; ii++ {
		s := <-config.ReportChan
		if (s.DataNum == 0) && (s.FileNum == 0) {
			//logs.Log.App(" *     [任务小计：%s | KEYIN：%s]   无采集结果，用时 %v！\n", s.SpiderName, s.Keyin, s.Time)
			continue
		}
		//logs.Log.Informational(" * ")
		//switch {
		//case s.DataNum > 0 && s.FileNum == 0:
		//	logs.Log.App(" *     [任务小计：%s | KEYIN：%s]   共采集数据 %v 条，用时 %v！\n",
		//		s.SpiderName, s.Keyin, s.DataNum, s.Time)
		//case s.DataNum == 0 && s.FileNum > 0:
		//	logs.Log.App(" *     [任务小计：%s | KEYIN：%s]   共下载文件 %v 个，用时 %v！\n",
		//		s.SpiderName, s.Keyin, s.FileNum, s.Time)
		//default:
		//	logs.Log.App(" *     [任务小计：%s | KEYIN：%s]   共采集数据 %v 条 + 下载文件 %v 个，用时 %v！\n",
		//		s.SpiderName, s.Keyin, s.DataNum, s.FileNum, s.Time)
		//}

		self.sum[0] += s.DataNum
		self.sum[1] += s.FileNum
	}

	// 总耗时
	self.takeTime = time.Since(config.StartTime)
	var prefix = func() string {
		if self.Status() == config.STOP {
			return "任务中途取消："
		}
		return "本次"
	}()
	// 打印总结报告
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(` *********************************************************************************************************************************** `)
	//logs.Log.Informational(" * ")
	switch {
	case self.sum[0] > 0 && self.sum[1] == 0:
		//logs.Log.App(" *                            —— %s合计采集【数据 %v 条】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
		//	prefix, self.sum[0], config.GetPageCount(1), config.GetPageCount(-1), config.GetPageCount(0), self.takeTime)
		fmt.Printf(" *                            —— %s合计采集【数据 %v 条】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
			prefix, self.sum[0], config.GetPageCount(1), config.GetPageCount(-1), config.GetPageCount(0), self.takeTime)
	case self.sum[0] == 0 && self.sum[1] > 0:
		//logs.Log.App(" *                            —— %s合计采集【文件 %v 个】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
		//	prefix, self.sum[1], config.GetPageCount(1), config.GetPageCount(-1), config.GetPageCount(0), self.takeTime)
		fmt.Printf(" *                            —— %s合计采集【文件 %v 个】， 实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
			prefix, self.sum[1], config.GetPageCount(1), config.GetPageCount(-1), config.GetPageCount(0), self.takeTime)
	case self.sum[0] == 0 && self.sum[1] == 0:
		fmt.Printf(" *                            —— %s无采集结果，实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
			prefix, config.GetPageCount(1), config.GetPageCount(-1), config.GetPageCount(0), self.takeTime)
	default:
		fmt.Printf(" *                            —— %s合计采集【数据 %v 条 + 文件 %v 个】，实爬【成功 %v URL + 失败 %v URL = 合计 %v URL】，耗时【%v】 ——",
			prefix, self.sum[0], self.sum[1], config.GetPageCount(1), config.GetPageCount(-1), config.GetPageCount(0), self.takeTime)
	}
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(` *********************************************************************************************************************************** `)

	// 单机模式并发运行，需要标记任务结束
	if self.SysConfig.Mode == config.OFFLINE {
		//self.LogRest()
		self.finishOnce.Do(func() { close(self.finish) })
	}
}

// 客户端向服务端反馈日志
//func (self *System) socketLog() {
//	for self.canSocketLog {
//		_, msg, ok := logs.Log.StealOne()
//		if !ok {
//			return
//		}
//		if self.Teleport.CountNodes() == 0 {
//			// 与服务器失去连接后，抛掉返馈日志
//			continue
//		}
//		self.Teleport.Request(msg, "log", "")
//	}
//}

func (self *System) checkPort() bool {
	if self.SysConfig.Port == 0 {
		//logs.Log.Warning(" *     —— 亲，分布式端口不能为空哦~")
		return false
	}
	return true
}

func (self *System) checkAll() bool {
	if self.SysConfig.Master == "" || !self.checkPort() {
		//logs.Log.Warning(" *     —— 亲，服务器地址不能为空哦~")
		return false
	}
	return true
}

// 设置任务运行时公共配置
func (self *System) setAppConf(task *mission.Task) {
	self.SysConfig.ThreadNum = task.ThreadNum
	self.SysConfig.Pausetime = task.Pausetime
	self.SysConfig.OutType = task.OutType
	self.SysConfig.DockerCap = task.DockerCap
	self.SysConfig.SuccessInherit = task.SuccessInherit
	self.SysConfig.FailureInherit = task.FailureInherit
	self.SysConfig.Limit = task.Limit
	self.SysConfig.ProxyMinute = task.ProxyMinute
	self.SysConfig.Keyins = task.Keyins
}
func (self *System) setTask(task *mission.Task) {
	task.ThreadNum = self.SysConfig.ThreadNum
	task.Pausetime = self.SysConfig.Pausetime
	task.OutType = self.SysConfig.OutType
	task.DockerCap = self.SysConfig.DockerCap
	task.SuccessInherit = self.SysConfig.SuccessInherit
	task.FailureInherit = self.SysConfig.FailureInherit
	task.Limit = self.SysConfig.Limit
	task.ProxyMinute = self.SysConfig.ProxyMinute
	task.Keyins = self.SysConfig.Keyins
}