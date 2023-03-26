package zlock

import (
	"sync"
)

var rewritelock sync.Mutex

type ZLockMap struct {
	SyncMap *sync.Map
}

func New() *ZLockMap {
	return &ZLockMap{SyncMap: &sync.Map{}}
}

func (zm *ZLockMap) Lock(key interface{}) {
	val, ok := zm.SyncMap.Load(key)
	if ok {
		l := val.(*sync.Mutex)
		l.Lock()
	}
	rewritelock.Lock()
	defer rewritelock.Unlock()
	val, ok = zm.SyncMap.Load(key)
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
