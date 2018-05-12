package spider

import (
	"fmt"
	"github.com/alezh/novel/system/pinyin"
)

type Species struct {
	list   []*Spider
	hash   map[string]*Spider
	sorted bool
}

var SpeciesCollection = &Species{
	list:make([]*Spider,0),
	hash:make(map[string]*Spider,0),
}

//载入规则
func (s *Species)Load(spider *Spider) *Spider {
	name := spider.Name
	for i := 2; true; i++ {
		if _, ok := s.hash[name]; !ok {
			s.hash[name] = spider
			break
		}
		name = fmt.Sprintf("%s(%d)", spider.Name, i)
	}
	fmt.Println(name)
	spider.Name = name
	s.list = append(s.list,spider)
	return spider
}

func (self *Species) Get () []*Spider {
	if !self.sorted {
		l := len(self.list)
		initials := make([]string, l)
		newlist := map[string]*Spider{}
		for i := 0; i < l; i++ {
			initials[i] = self.list[i].GetName()
			newlist[initials[i]] = self.list[i]
		}
		pinyin.SortInitials(initials)
		for i := 0; i < l; i++ {
			self.list[i] = newlist[initials[i]]
		}
		self.sorted = true
	}
	return self.list
}

func (self *Species) GetByName(name string) *Spider {
	return self.hash[name]
}