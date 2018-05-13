package data

import (
	"sort"
	"github.com/alezh/novel/system/data/collector"
	"github.com/alezh/novel/config"
)

// 初始化输出方式列表collector.DataOutputLib
func init() {
	for out, _ := range collector.DataOutput {
		collector.DataOutputLib = append(collector.DataOutputLib, out)
	}
	sort.Strings(collector.DataOutputLib)
}

// 刷新输出方式的状态
func RefreshOutput() {
	switch config.Task.OutType {
	case "mgo":
		//mgo.Refresh()
	case "mysql":
		//mysql.Refresh()
	//case "kafka":
		//kafka.Refresh()
	}
}
