package data

import "sync"

//输出数据pool
type (
	Data map[string]interface{}
)

var dataPool = &sync.Pool{
	New: func() interface{} {
		return Data{}
	},
}

func GetData(Data)  Data{
	data :=dataPool.Get().(Data)
	return data
}

func PutData(data Data)  {
	dataPool.Put(data)
}