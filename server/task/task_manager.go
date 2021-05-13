package task

import (
	"context"
	"crawshaw.io/sqlite/sqlitex"
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/web-bt-client/db"
	"log"
	"sync"
)

type Manager struct {
	client  *torrent.Client
	taskMap *Map
}

func (tm *Manager) taskExists(infoHash string) bool {
	if tm.taskMap.HasMember(infoHash) {
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

func (tm *Manager) newTorrentTask(t *torrent.Torrent) (*TorrentTask, error) {
	if tm.taskExists(t.InfoHash().String()) {
		return nil, fmt.Errorf("任务已存在 %s", t.InfoHash().String())
	} else {
		task := NewTorrentTask(t, tm.client)
		tm.taskMap.Store(t.InfoHash().String(), task)
		return task, nil
	}
}

func (tm *Manager) AddUriTask(uri string) (string, error) {
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

func (tm *Manager) Download(hash string, files []string) error {
	if task := tm.taskMap.Load(hash); task != nil {
		task.Download(files)
		return nil
	}
	return fmt.Errorf("任务 %s 不存在", hash)
}

func (tm *Manager) Stop(hash string) error {
	if task := tm.taskMap.Load(hash); task != nil {
		task.Stop()
		return nil
	}
	return fmt.Errorf("任务 %s 不存在", hash)
}

func (tm *Manager) Start(hash string) error {
	if task := tm.taskMap.Load(hash); task != nil {
		return task.Start()
	}
	return fmt.Errorf("任务 %s 不存在", hash)
}

func (tm *Manager) GetTorrentInfo(hash string) (TorrentInfoWrapper, error) {
	if task := tm.taskMap.Load(hash); task != nil {
		return task.GetTorrentInfo(), nil
	}
	// todo tasks内加载
	return TorrentInfoWrapper{}, fmt.Errorf("未找到 %s 信息", hash)
}

var taskManager *Manager
var tmOnce sync.Once

func GetTaskManager() *Manager {
	tmOnce.Do(func() {
		cfg := torrent.NewDefaultClientConfig()
		cfg.Seed = true
		client, err := torrent.NewClient(cfg)
		if err != nil {
			log.Fatalln(err)
		}
		taskManager = &Manager{
			client:  client,
			taskMap: NewTaskMap(),
		}
	})
	return taskManager
}
