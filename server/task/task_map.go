package task

import (
	"sync"
	"sync/atomic"
)

type Map struct {
	count int32
	cMap  sync.Map
}

func (tm *Map) Count() int32 {
	return atomic.LoadInt32(&tm.count)
}

func (tm *Map) HasMember(hash string) bool {
	if _, ok := tm.cMap.Load(hash); ok {
		return true
	}
	return false
}

func (tm *Map) Load(hash string) *TorrentTask {
	if v, ok := tm.cMap.Load(hash); ok {
		return v.(*TorrentTask)
	}
	return nil
}

func (tm *Map) Store(hash string, task *TorrentTask) {
	tm.cMap.Store(hash, task)
	atomic.AddInt32(&tm.count, 1)
}

func (tm *Map) Delete(hash string) {
	if _, ok := tm.cMap.Load(hash); ok {
		tm.cMap.Delete(hash)
		atomic.AddInt32(&tm.count, -1)
	}
}

func NewTaskMap() *Map {
	return &Map{
		count: 0,
	}
}
