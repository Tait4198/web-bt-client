package ws

import (
	"sync"
	"sync/atomic"
)

type ClientMap struct {
	count int32
	cMap  sync.Map
}

func (cm *ClientMap) Count() int32 {
	return atomic.LoadInt32(&cm.count)
}

func (cm *ClientMap) HasMember(hash string) bool {
	if _, ok := cm.cMap.Load(hash); ok {
		return true
	}
	return false
}

func (cm *ClientMap) Load(id string) *WebSocketClient {
	if v, ok := cm.cMap.Load(id); ok {
		return v.(*WebSocketClient)
	}
	return nil
}

func (cm *ClientMap) Store(id string, wsc *WebSocketClient) {
	cm.cMap.Store(id, wsc)
	atomic.AddInt32(&cm.count, 1)
}

func (cm *ClientMap) Delete(id string) {
	if _, ok := cm.cMap.Load(id); ok {
		cm.cMap.Delete(id)
		atomic.AddInt32(&cm.count, -1)
	}
}

func (cm *ClientMap) Range(f func(id string, wsc *WebSocketClient)) {
	cm.cMap.Range(func(key, value interface{}) bool {
		id := key.(string)
		wsc := value.(*WebSocketClient)
		f(id, wsc)
		return true
	})
}

func NewClientMap() *ClientMap {
	return &ClientMap{
		count: 0,
	}
}
