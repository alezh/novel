package reptilian

import (
	. "github.com/alezh/novel/system/spider"
	"github.com/henrylee2cn/pholcus/common/util"
)

type (

	SpiderQueue interface {
		Reset() //重置清空队列
		Add(*Spider)
		AddAll([]*Spider)
		AddKeyins(string) //为队列成员遍历添加Keyin属性，但前提必须是队列成员未被添加过keyin
		GetByIndex(int) *Spider
		GetByName(string) *Spider
		GetAll() []*Spider
		Len() int // 返回队列长度
	}
	rq struct {
		list []*Spider
	}
)

func NewPool() SpiderQueue {
	return &rq{
		list: []*Spider{},
	}
}

func (self *rq) Reset() {
	self.list = []*Spider{}
}

func (self *rq) Add(sp *Spider) {
	sp.SetId(self.Len())
	self.list = append(self.list, sp)
}

func (self *rq) AddAll(list []*Spider) {
	for _, v := range list {
		self.Add(v)
	}
}

// 添加keyin，遍历蜘蛛队列得到新的队列（已被显式赋值过的spider将不再重新分配Keyin）
func (self *rq) AddKeyins(keyins string) {
	keyinSlice := util.KeyinsParse(keyins)
	if len(keyinSlice) == 0 {
		return
	}

	unit1 := []*Spider{} // 不可被添加自定义配置的蜘蛛
	unit2 := []*Spider{} // 可被添加自定义配置的蜘蛛
	for _, v := range self.GetAll() {
		if v.GetKeyin() == KEYIN {
			unit2 = append(unit2, v)
			continue
		}
		unit1 = append(unit1, v)
	}

	if len(unit2) == 0 {
		//logs.Log.Warning("本批任务无需填写自定义配置！\n")
		return
	}

	self.Reset()

	for _, keyin := range keyinSlice {
		for _, v := range unit2 {
			v.Keyin = keyin
			nv := *v
			self.Add((&nv).Copy())
		}
	}
	if self.Len() == 0 {
		self.AddAll(append(unit1, unit2...))
	}

	self.AddAll(unit1)
}

func (self *rq) GetByIndex(idx int) *Spider {
	return self.list[idx]
}

func (self *rq) GetByName(n string) *Spider {
	for _, sp := range self.list {
		if sp.GetName() == n {
			return sp
		}
	}
	return nil
}

func (self *rq) GetAll() []*Spider {
	return self.list
}

func (self *rq) Len() int {
	return len(self.list)
}