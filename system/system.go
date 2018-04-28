package system

import (
	"github.com/henrylee2cn/teleport"
	"time"
	"sync"
	_ "github.com/alezh/novel/lib"
	"github.com/alezh/novel/system/spider"
	"github.com/alezh/novel/config"
	"io"
)

type (
	Enging interface {

	}
	System struct {
		*config.SysConfig                   // 全局配置
		*spider.Species                     // 全部蜘蛛种类
		teleport.Teleport                   // socket长连接双工通信接口，json数据传输
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
		Teleport:    teleport.New(),

	}
	return sys
}

func (sys *System)Init(mode int, port int, master string, w ...io.Writer)  {
	sys.canSocketLog = false
	if len(w) > 0 {
		//sys.SetLog(w[0])
	}
	//sys.LogGoOn()
	sys.SysConfig.Mode, sys.SysConfig.Port, sys.SysConfig.Master = mode, port, master
	sys.Teleport = teleport.New()

}

func (sys *System)ReInit()  {

}

func (sys *System)Start()  {

}

func (sys *System)Stop()  {

}

func (sys *System)GetConfig()  {
	
}

func (sys *System)SetConfig()  {
	
}

func (sys *System)LogGoOn()  {
	
}


