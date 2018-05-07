package reptilian

import (
	"sync"
	"github.com/alezh/novel/config"
)

type (
	ReptilianPool interface {

	}
	rp struct {
		capacity int
		count    int
		usable   chan Reptilian
		all      []Reptilian
		status   int
		sync.RWMutex
	}
)

func NewReptilianPool() ReptilianPool {
	return &rp{
		status: config.RUN,
		all:    make([]Reptilian, 0, 50), // 蜘蛛池最大容量
	}
}