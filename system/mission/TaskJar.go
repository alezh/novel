package mission

import (
	"github.com/henrylee2cn/teleport"
	"encoding/json"
)

type (
	// 分布式的接口
	Distributer interface {
		// 主节点从仓库发送一个任务
		Send(clientNum int) Task
		// 从节点接收一个任务到仓库
		Receive(task *Task)
		// 返回与之连接的节点数
		CountNodes() int
	}
	// 用于分布式分发的任务
	Task struct {
		Id             int
		Spiders        []map[string]string // 蜘蛛规则name字段与keyin字段，规定格式map[string]string{"name":"baidu","keyin":"henry"}
		ThreadNum      int                 // 全局最大并发量
		Pausetime      int64               // 暂停时长参考/ms(随机: Pausetime/2 ~ Pausetime*2)
		OutType        string              // 输出方式
		DockerCap      int                 // 分段转储容器容量
		DockerQueueCap int                 // 分段输出池容量，不小于2
		SuccessInherit bool                // 继承历史成功记录
		FailureInherit bool                // 继承历史失败记录
		Limit          int64               // 采集上限，0为不限，若在规则中设置初始值为LIMIT则为自定义限制，否则默认限制请求数
		ProxyMinute    int64               // 代理IP更换的间隔分钟数
		// 选填项
		Keyins string // 自定义输入，后期切分为多个任务的Keyin自定义配置
	}
	// 任务仓库
	TaskJar struct {
		Tasks chan *Task
	}
	masterTaskHandle struct {
		Distributer
	}
	// 主节点自动接收从节点消息并打印的操作
	masterLogHandle struct{}
)

func NewTaskJar() *TaskJar {
	return &TaskJar{
		Tasks: make(chan *Task, 1024),
	}
}

// 服务器向仓库添加一个任务
func (self *TaskJar) Push(task *Task) {
	id := len(self.Tasks)
	task.Id = id
	self.Tasks <- task
}

// 客户端从本地仓库获取一个任务
func (self *TaskJar) Pull() *Task {
	return <-self.Tasks
}

// 仓库任务总数
func (self *TaskJar) Len() int {
	return len(self.Tasks)
}

// 主节点从仓库发送一个任务
func (self *TaskJar) Send(clientNum int) Task {
	return *<-self.Tasks
}

// 从节点接收一个任务到仓库
func (self *TaskJar) Receive(task *Task) {
	self.Tasks <- task
}


// 创建主节点API
func MasterApi(n Distributer) teleport.API {
	return teleport.API{
		// 分配任务给客户端
		"task": &masterTaskHandle{n},

		// 打印接收到的日志
		"log": &masterLogHandle{},
	}
}

func (self *masterTaskHandle) Process(receive *teleport.NetData) *teleport.NetData {
	b, _ := json.Marshal(self.Send(self.CountNodes()))
	return teleport.ReturnData(string(b))
}

func (*masterLogHandle) Process(receive *teleport.NetData) *teleport.NetData {
	//logs.Log.Informational(" * ")
	//logs.Log.Informational(" *     [ %s ]    %s", receive.From, receive.Body)
	//logs.Log.Informational(" * ")
	return nil
}