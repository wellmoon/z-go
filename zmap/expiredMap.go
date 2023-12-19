package zmap

import (
	"fmt"
	"sync"
	"time"
)

type ExpiredMap struct {
	valueMap     map[interface{}]interface{}
	timeoutMap   map[interface{}]int64
	checkExpired bool
	lock         *sync.Mutex
}

func NewExpiredMap() *ExpiredMap {
	var lock sync.Mutex
	m := &ExpiredMap{make(map[interface{}]interface{}), make(map[interface{}]int64), false, &lock}
	return m
}

func (m *ExpiredMap) startCheck() {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.checkExpired {
		return
	} else {
		m.checkExpired = true
	}
	m.checkExpired = true
	go func() {
		for {
			if len(m.timeoutMap) == 0 {
				m.checkExpired = false
				break
			}
			for k, v := range m.timeoutMap {
				now := time.Now().UnixMilli()
				if now > v {
					delete(m.valueMap, k)
					delete(m.timeoutMap, k)
				}
			}
			time.Sleep(time.Duration(500) * time.Millisecond)
		}
	}()
}

func (m *ExpiredMap) Put(key interface{}, val interface{}, timeoutSecond int) {
	m.startCheck()
	m.lock.Lock()
	defer m.lock.Unlock()
	m.timeoutMap[key] = time.Now().UnixMilli() + int64(timeoutSecond)*1000
	m.valueMap[key] = val
}

func (m *ExpiredMap) GetVal(key interface{}) interface{} {
	m.startCheck()
	m.lock.Lock()
	return m.valueMap[key]
}

func (m *ExpiredMap) Contains(key interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.valueMap[key]
	return ok
}

func (m *ExpiredMap) String() string {
	m.startCheck()
	m.lock.Lock()
	res := ""
	for k, v := range m.valueMap {
		res = fmt.Sprintf("%s%v:%v %d\n", res, k, v, m.timeoutMap[k])
	}
	return res
}

func (m *ExpiredMap) IsNil() bool {
	m.startCheck()
	m.lock.Lock()
	return len(m.timeoutMap) == 0
}
