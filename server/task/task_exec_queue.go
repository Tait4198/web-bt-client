package task

import (
	"container/list"
	"fmt"
)

// todo 优先下载

type ExecQueue struct {
	tasks  *list.List
	waitCh chan byte
	execCh chan byte
}

func (q *ExecQueue) pushBack(task *TorrentTask) {
	q.tasks.PushBack(task)
	q.waitCh <- 0
}

func (q *ExecQueue) find(hash string) (*list.Element, int) {
	index := 0
	for e := q.tasks.Front(); e != nil; e = e.Next() {
		index++
		fmt.Printf("%s %s\n", hash, e.Value.(*TorrentTask).param.InfoHash)
		if e.Value.(*TorrentTask).param.InfoHash == hash {
			return e, index
		}
	}
	return nil, -1
}

func (q *ExecQueue) run() {
	for range q.waitCh {
		q.execCh <- 0
		if e := q.tasks.Front(); e != nil {
			q.tasks.Remove(e)
			task := e.Value.(*TorrentTask)
			BroadcastTaskStatus(task, QueueStatus, false)
			go func() {
				BroadcastTaskStatus(task, Pause, false)
				defer func() {
					<-q.execCh
				}()
				task.exec()
			}()
		} else {
			<-q.execCh
		}
	}
}

func (q *ExecQueue) remove(hash string) {
	if e, _ := q.find(hash); e != nil {
		q.tasks.Remove(e)
		BroadcastTaskStatus(e.Value.(*TorrentTask), QueueStatus, false)
	}
}

func NewExecQueue(maxExecSize int) *ExecQueue {
	q := &ExecQueue{
		tasks:  list.New(),
		waitCh: make(chan byte, 1<<12),
		execCh: make(chan byte, maxExecSize),
	}
	go q.run()
	return q
}
