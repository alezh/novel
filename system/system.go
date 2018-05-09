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
)

type (
	Enging interface {
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
		SpiderQueue: reptilian.NewPool(),
		Teleport:    teleport.New(),
		TaskJar:     mission.NewTaskJar(),
		ReptilianPool:reptilian.NewPool(),
	}
	return sys
}
//
//func WebStart()  {
//	app.NewApp().Run(iris.Addr(":8080"),app.Configure)
//}

func (sys *System)Init(mode int, port int, master string, w ...io.Writer)  {
	sys.canSocketLog = false
	if len(w) > 0 {
		//sys.SetLog(w[0])
	}
	//sys.LogGoOn()
	sys.SysConfig.Mode, sys.SysConfig.Port, sys.SysConfig.Master = mode, port, master
	sys.Teleport = teleport.New()
	sys.TaskJar = mission.NewTaskJar()
	sys.SpiderQueue = reptilian.NewPool()
}

func (sys *System)ReInit()  {

}

func (sys *System)Start()  {

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


