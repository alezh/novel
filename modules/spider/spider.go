package spider

import "github.com/alezh/novel/modules/rule"

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
	RuleTree        *rule.RuleTree                                             // 定义具体的采集规则树

}

func (s *Spider)Register()  *Spider{
	return SpeciesCollection.Load(s)
}

func (s *Spider)GetAll() []*Spider {
	return SpeciesCollection.Get()
}