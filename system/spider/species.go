package spider

import "fmt"

type Species struct {
	list   []*Spider
	hash   map[string]*Spider
}

var SpeciesCollection = &Species{
	make([]*Spider,0),
	make(map[string]*Spider,0),
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

func (s *Species) Get () []*Spider {
	if len(s.list)>0{
		return s.list
	}
	return nil
}