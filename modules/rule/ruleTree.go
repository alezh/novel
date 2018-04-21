package rule

import "github.com/alezh/novel/system"

type RuleTree struct {
	Root  func(*system.Mission)   // 根节点(执行入口)
	Trunk map[string]*Rule // 节点散列表(执行采集过程)
}

type Rule struct {
	ItemFields []string                                                  // 结果字段列表(选填，写上可保证字段顺序)
	ParseFunc  func(*system.Mission)                                     // 内容解析函数
	AidFunc    func(*system.Mission,map[string]interface{}) interface{}  // 通用辅助函数
}