package data

import (
	"sync"
)

type (
	// 数据存储单元
	Data map[string]interface{}
)

var (
	dataCellPool = &sync.Pool{
		New: func() interface{} {
			return Data{}
		},
	}
)

func GetDataCell(ruleName string, data map[string]interface{}, url string, parentUrl string, downloadTime string) Data {
	cell := dataCellPool.Get().(Data)
	cell["RuleName"] = ruleName   //规定Data中的key
	cell["Data"] = data           //数据存储,key须与Rule的Fields保持一致
	cell["Url"] = url             //用于索引
	cell["ParentUrl"] = parentUrl //DataCell的上级url
	cell["DownloadTime"] = downloadTime
	return cell
}


func PutDataCell(cell Data) {
	cell["RuleName"] = nil
	cell["Data"] = nil
	cell["Url"] = nil
	cell["ParentUrl"] = nil
	cell["DownloadTime"] = nil
	dataCellPool.Put(cell)
}
