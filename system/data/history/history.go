package history

import "sync"

type History struct {
	*Success
	*Failure
	provider string
	sync.RWMutex
}