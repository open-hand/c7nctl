package c7nclient

import "sync"

type SafeMap struct {
	sync.RWMutex
	Map map[string]interface{}
}

func NewSafeMap() *SafeMap {
	sm := new(SafeMap)
	sm.Map = make(map[string]interface{})
	return sm

}

func (sm *SafeMap) readMap(key string) interface{} {
	sm.RLock()
	value := sm.Map[key]
	sm.RUnlock()
	return value
}

func (sm *SafeMap) writeMap(key string, value interface{}) {
	sm.Lock()
	sm.Map[key] = value
	sm.Unlock()
}
