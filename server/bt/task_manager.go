package bt

import (
	"context"
	"crawshaw.io/sqlite/sqlitex"
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/web-bt-client/db"
	"log"
	"sync"
)

type TaskManager struct {
	client  *torrent.Client
	taskMap sync.Map
}

func (tm *TaskManager) taskExists(infoHash string) bool {
	if _, ok := tm.taskMap.Load(infoHash); ok {
		return true
	}
	conn := db.GetPool().Get(context.TODO())
	defer db.GetPool().Put(conn)

	stmt := conn.Prep("select count(*) from tasks where info_hash = $hash")
	stmt.SetText("$hash", infoHash)
	if v, err := sqlitex.ResultInt(stmt); err == nil {
		return v > 0
	} else {
		log.Printf("Query Count Error %s", err)
	}
	return true
}

func (tm *TaskManager) newTorrentTask(t *torrent.Torrent) (*TorrentTask, error) {
	if tm.taskExists(t.InfoHash().String()) {
		return nil, fmt.Errorf("任务已存在 %s", t.InfoHash().String())
	} else {
		task := NewTorrentTask(t, tm.client)
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
		//task.GetInfo()
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
