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
		//task.Download()
		task.GetInfo()
		return t.InfoHash().String(), nil
	}
	return "", err
}

func (tm *TaskManager) Download(hash string, files []string) error {
	if task, ok := tm.taskMap.Load(hash); ok {
		(task.(*TorrentTask)).Download(files)
		return nil
	}
	return fmt.Errorf("任务 %s 不存在", hash)
}

func (tm *TaskManager) Stop(hash string) error {
	if task, ok := tm.taskMap.Load(hash); ok {
		(task.(*TorrentTask)).Stop()
	}
	return fmt.Errorf("任务 %s 不存在", hash)
}

func (tm *TaskManager) GetTorrentInfo(hash string) (TaskTorrentInfo, error) {
	if task, ok := tm.taskMap.Load(hash); ok {
		return (task.(*TorrentTask)).GetTaskTorrentInfo(), nil
	}
	// todo tasks内加载
	return TaskTorrentInfo{}, fmt.Errorf("未找到 %s 信息", hash)
}

var taskManager *TaskManager
var tmOnce sync.Once

func GetTaskManager() *TaskManager {
	tmOnce.Do(func() {
		cfg := torrent.NewDefaultClientConfig()
		cfg.Seed = true
		client, err := torrent.NewClient(cfg)
		if err != nil {
			log.Fatalln(err)
		}
		taskManager = &TaskManager{
			client: client,
		}
	})
	return taskManager
}
