package types

import (
	"math/rand"
	"sync"
	"time"
)

type UniqueRandomIDSource struct {
	l    uint64
	has  map[string]bool
	r    *rand.Rand
	lock sync.Mutex
}

func NewUniqueRandomIDSource(l uint64) *UniqueRandomIDSource {
	ret := &UniqueRandomIDSource{
		l:   l,
		has: make(map[string]bool),
		r:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return ret
}

func (this *UniqueRandomIDSource) Get() string {
	this.lock.Lock()
	defer this.lock.Unlock()
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for {
		bytes := []byte(str)
		result := []byte{}
		for i := uint64(0); i < this.l; i++ {
			result = append(result, bytes[this.r.Intn(len(bytes))])
		}
		if !this.has[string(result)] {
			this.has[string(result)] = true
			return string(result)
		}
	}

}
