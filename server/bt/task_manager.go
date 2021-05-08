package bt

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"log"
	"sync"
)

type TaskManager struct {
	client  *torrent.Client
	taskMap sync.Map
}

func (tm *TaskManager) newTorrentTask(t *torrent.Torrent) (*TorrentTask, error) {
	if _, ok := tm.taskMap.Load(t.InfoHash().String()); ok {
		return nil, fmt.Errorf("任务已存在 %s", t.InfoHash().String())
	} else {
		task := NewTorrentTask(t)
		tm.taskMap.Store(t.InfoHash().String(), task)
		return task, nil
	}
}

func (tm *TaskManager) AddUriTask(uri string) (string, error) {
	t, err := tm.client.AddMagnet(uri)
	if err != nil {
		return "", err
	}
	task, err := tm.newTorrentTask(t)
	if err == nil {
		task.Download()
		return t.InfoHash().String(), nil
	}
	return "", err
}

func (tm *TaskManager) StopTask(hash string) {
	if task, ok := tm.taskMap.Load(hash); ok {
		(task.(*TorrentTask)).Stop()
	}
}

var taskManager *TaskManager
var tmOnce sync.Once

func GetTaskManager() *TaskManager {
	tmOnce.Do(func() {
		client, err := torrent.NewClient(nil)
		if err != nil {
			log.Fatalln(err)
		}
		taskManager = &TaskManager{
			client: client,
		}
	})
	return taskManager
}
