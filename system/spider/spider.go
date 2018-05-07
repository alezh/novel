package spider

import (
	"sync"
	"github.com/alezh/novel/system/mission"
	"github.com/henrylee2cn/pholcus/common/util"
	"math"
	"github.com/alezh/novel/config"
	"github.com/alezh/novel/system/http/request"
	"time"
)

const (
	KEYIN       = util.USE_KEYIN // 若使用Spider.Keyin，则须在规则中设置初始值为USE_KEYIN
	LIMIT       = math.MaxInt64  // 如希望在规则中自定义控制Limit，则Limit初始值必须为LIMIT
	FORCED_STOP = "——主动终止Spider——"
)

//规则添加 获取规则相关数据

type Spider struct {
	Name            string                                                     // 用户界面显示的名称（应保证唯一性）
	Description     string                                                     // 用户界面显示的描述
	//Pausetime       int64                                                      // 随机暂停区间(50%~200%)，若规则中直接定义，则不被界面传参覆盖
	Limit           int64                                                      // 默认限制请求数，0为不限；若规则中定义为LIMIT，则采用规则的自定义限制方案
	Keyin           string                                                     // 自定义输入的配置信息，使用前须在规则中设置初始值为KEYIN
	EnableCookie    bool                                                       // 所有请求是否使用cookie记录
	//NotDefaultField bool                                                       // 是否禁止输出结果中的默认字段 Url/ParentUrl/DownloadTime
	Namespace       func(self *Spider) string                                  // 命名空间，用于输出文件、路径的命名
	SubNamespace    func(self *Spider, dataCell map[string]interface{}) string // 次级命名，用于输出文件、路径的命名，可依赖具体数据内容
	RuleTree        *RuleTree                                             // 定义具体的采集规则树

	// 以下字段系统自动赋值
	id        int               // 自动分配的SpiderQueue中的索引
	subName   string            // 由Keyin转换为的二级标识名
	reqMatrix *mission.Matrix   // 请求矩阵
	timer     *Timer          // 定时器
	status    int               // 执行状态
	lock      sync.RWMutex
	once      sync.Once

}

func (s *Spider)Register()  *Spider{
	return SpeciesCollection.Load(s)
}

// 指定规则的获取结果的字段名列表
func (self *Spider) GetItemFields(rule *Rule) []string {
	return rule.ItemFields
}

// 返回结果字段名的值
// 不存在时返回空字符串
func (self *Spider) GetItemField(rule *Rule, index int) (field string) {
	if index > len(rule.ItemFields)-1 || index < 0 {
		return ""
	}
	return rule.ItemFields[index]
}

// 返回结果字段名的其索引
// 不存在时索引为-1
func (self *Spider) GetItemFieldIndex(rule *Rule, field string) (index int) {
	for idx, v := range rule.ItemFields {
		if v == field {
			return idx
		}
	}
	return -1
}

// 为指定Rule动态追加结果字段名，并返回索引位置
// 已存在时返回原来索引位置
func (self *Spider) UpsertItemField(rule *Rule, field string) (index int) {
	for i, v := range rule.ItemFields {
		if v == field {
			return i
		}
	}
	rule.ItemFields = append(rule.ItemFields, field)
	return len(rule.ItemFields) - 1
}



func (s *Spider)GetAll() []*Spider {
	return SpeciesCollection.Get()
}

// 返回是否作为新的失败请求被添加至队列尾部
func (self *Spider) DoHistory(req *request.Request, ok bool) bool {
	return self.reqMatrix.DoHistory(req, ok)
}

func (self *Spider) RequestPush(req *request.Request) {
	self.reqMatrix.Push(req)
}

func (self *Spider) RequestPull() *request.Request {
	return self.reqMatrix.Pull()
}

func (self *Spider) RequestUse() {
	self.reqMatrix.Use()
}

func (self *Spider) RequestFree() {
	self.reqMatrix.Free()
}

func (self *Spider) RequestLen() int {
	return self.reqMatrix.Len()
}

func (self *Spider) TryFlushSuccess() {
	self.reqMatrix.TryFlushSuccess()
}

func (self *Spider) TryFlushFailure() {
	self.reqMatrix.TryFlushFailure()
}


// 获取蜘蛛名称
func (self *Spider) GetName() string {
	return self.Name
}

// 获取蜘蛛二级标识名
func (self *Spider) GetSubName() string {
	self.once.Do(func() {
		self.subName = self.GetKeyin()
		if len([]rune(self.subName)) > 8 {
			self.subName = util.MakeHash(self.subName)
		}
	})
	return self.subName
}

// 安全返回指定规则
func (self *Spider) GetRule(ruleName string) (*Rule, bool) {
	rule, found := self.RuleTree.Trunk[ruleName]
	return rule, found
}

// 返回指定规则
func (self *Spider) MustGetRule(ruleName string) *Rule {
	return self.RuleTree.Trunk[ruleName]
}

// 返回规则树
func (self *Spider) GetRules() map[string]*Rule {
	return self.RuleTree.Trunk
}

// 获取蜘蛛描述
func (self *Spider) GetDescription() string {
	return self.Description
}
// 获取蜘蛛ID
func (self *Spider) GetId() int {
	return self.id
}

// 设置蜘蛛ID
func (self *Spider) SetId(id int) {
	self.id = id
}
// 获取自定义配置信息
func (self *Spider) GetKeyin() string {
	return self.Keyin
}

// 设置自定义配置信息
func (self *Spider) SetKeyin(keyword string) {
	self.Keyin = keyword
}

// 获取采集上限
// <0 表示采用限制请求数的方案
// >0 表示采用规则中的自定义限制方案
func (self *Spider) GetLimit() int64 {
	return self.Limit
}

// 设置采集上限
// <0 表示采用限制请求数的方案
// >0 表示采用规则中的自定义限制方案
func (self *Spider) SetLimit(max int64) {
	self.Limit = max
}

// 控制所有请求是否使用cookie
func (self *Spider) GetEnableCookie() bool {
	return self.EnableCookie
}

// 设置定时器
// @id为定时器唯一标识
// @bell==nil时为倒计时器，此时@tol为睡眠时长
// @bell!=nil时为闹铃，此时@tol用于指定醒来时刻（从now起遇到的第tol个bell）
func (self *Spider) SetTimer(id string, tol time.Duration, bell *Bell) bool {
	if self.timer == nil {
		self.timer = newTimer()
	}
	return self.timer.set(id, tol, bell)
}

// 启动定时器，并返回定时器是否可以继续使用
func (self *Spider) RunTimer(id string) bool {
	if self.timer == nil {
		return false
	}
	return self.timer.sleep(id)
}

// 返回一个自身复制品
func (self *Spider) Copy() *Spider {
	ghost := &Spider{}
	ghost.Name = self.Name
	ghost.subName = self.subName

	ghost.RuleTree = &RuleTree{
		Root:  self.RuleTree.Root,
		Trunk: make(map[string]*Rule, len(self.RuleTree.Trunk)),
	}
	for k, v := range self.RuleTree.Trunk {
		ghost.RuleTree.Trunk[k] = new(Rule)

		ghost.RuleTree.Trunk[k].ItemFields = make([]string, len(v.ItemFields))
		copy(ghost.RuleTree.Trunk[k].ItemFields, v.ItemFields)

		ghost.RuleTree.Trunk[k].ParseFunc = v.ParseFunc
		ghost.RuleTree.Trunk[k].AidFunc = v.AidFunc
	}

	ghost.Description = self.Description
	//ghost.Pausetime = self.Pausetime
	ghost.EnableCookie = self.EnableCookie
	ghost.Limit = self.Limit
	ghost.Keyin = self.Keyin

	//ghost.NotDefaultField = self.NotDefaultField
	ghost.Namespace = self.Namespace
	ghost.SubNamespace = self.SubNamespace

	//ghost.timer = self.timer
	ghost.status = self.status

	return ghost
}

func (self *Spider) IsStopping() bool {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.status == config.STOP
}
// 若已主动终止任务，则崩溃爬虫协程
func (self *Spider) tryPanic() {
	if self.IsStopping() {
		panic(FORCED_STOP)
	}
}
