package zlock

import (
	"sync"
	"sync/atomic"
)

type ZLockMap struct {
	SyncMap *sync.Map
}

var opts int64 = 0
var innerMap = make(map[int64]*sync.Map, 0)
var lock sync.Mutex

func New() *ZLockMap {
	atomic.AddInt64(&opts, 1)
	return &ZLockMap{newMap()}
}

func newMap() *sync.Map {
	lock.Lock()
	defer lock.Unlock()
	var m sync.Map
	innerMap[opts] = &m
	return &m
}

func (zm *ZLockMap) Lock(key interface{}) {
	val, ok := zm.SyncMap.Load(key)
	if ok {
		l := val.(*sync.Mutex)
		l.Lock()
	} else {
		l := &sync.Mutex{}
		l.Lock()
		zm.SyncMap.Store(key, l)
	}
}

func (zm *ZLockMap) Unlock(key interface{}) {
	val, ok := zm.SyncMap.Load(key)
	if ok {
		l := val.(*sync.Mutex)
		l.Unlock()
		zm.SyncMap.Delete(key)
	}
}
