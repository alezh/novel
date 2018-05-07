package history

import (
	"sync"
	"github.com/alezh/novel/system/http/request"
)


type (
	Historier interface {

	}
	History struct {
		*Success
		*Failure
		provider string
		sync.RWMutex
	}
)

func New(name string, subName string) Historier{
	return &History{
		Success: &Success{
			//tabName:  util.FileNameReplace(successTabName),
			//fileName: successFileName,
			new:      make(map[string]bool),
			old:      make(map[string]bool),
		},
		Failure: &Failure{
			//tabName:  util.FileNameReplace(failureTabName),
			//fileName: failureFileName,
			list:     make(map[string]*request.Request),
		},
	}
}