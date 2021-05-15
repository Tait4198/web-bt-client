package task

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/web-bt-client/db"
	"log"
	"sync"
	"time"
)

type Manager struct {
	client  *torrent.Client
	taskMap *Map
}

func (tm *Manager) taskExists(infoHash string) bool {
	if tm.taskMap.HasMember(infoHash) {
		return true
	}
	return db.SelectTaskCount(infoHash) > 0
}

// 通过MetaInfo创建指定下载路径Torrent
func (tm *Manager) newMetaInfoTorrentWithPath(mi *metainfo.MetaInfo, path string) (*torrent.Torrent, error) {
	if path == "" {
		return tm.client.AddTorrent(mi)
	}
	spec := torrent.TorrentSpecFromMetaInfo(mi)
	spec.Storage = storage.NewMMap(path)
	if t, _, err := tm.client.AddTorrentSpec(spec); err == nil {
		return t, nil
	} else {
		return nil, err
	}
}

// 通过磁力链接创建指定下载路径Torrent
func (tm *Manager) newUriTorrentWithPath(uri string, path string) (*torrent.Torrent, error) {
	if path == "" {
		return tm.client.AddMagnet(uri)
	}
	spec, err := torrent.TorrentSpecFromMagnetUri(uri)
	if err != nil {
		return nil, err
	}
	spec.Storage = storage.NewMMap(path)
	if t, _, err := tm.client.AddTorrentSpec(spec); err == nil {
		return t, nil
	} else {
		return nil, err
	}
}

func (tm *Manager) recoveryTask(dbTask *db.Task, mi *metainfo.MetaInfo) error {
	if t, err := tm.newMetaInfoTorrentWithPath(mi, dbTask.DownloadPath); err == nil {
		infoHash := t.InfoHash().String()
		// 恢复下载进度
		if t.BytesCompleted() != dbTask.CompleteFileLength {
			if err := db.UpdateTaskCompleteFileLength(infoHash, t.BytesCompleted()); err != nil {
				return fmt.Errorf("任务 %s 下载进度恢复失败", infoHash)
			}
		}

		task := newTaskSetPath(t, tm, dbTask.DownloadPath)
		tm.taskMap.Store(t.InfoHash().String(), task)
		// todo 细化信息获取与下载
		go tm.TaskRun(task)
		return nil
	} else {
		return err
	}
}

func (tm *Manager) createTask(t *torrent.Torrent) (*TorrentTask, error) {
	if tm.taskExists(t.InfoHash().String()) {
		return nil, fmt.Errorf("任务已存在 %s", t.InfoHash().String())
	} else {
		dbTask := db.Task{
			InfoHash:     t.InfoHash().String(),
			DownloadPath: "D:\\Torrent",
			CreateTime:   time.Now().UnixNano(),
		}
		if t.Info() == nil {
			dbTask.TorrentName = t.InfoHash().String()
		} else {
			dbTask.TorrentName = t.Info().Name
		}
		if err := db.InsertTask(dbTask); err == nil {
			task := newTask(t, tm)
			tm.taskMap.Store(t.InfoHash().String(), task)
			go tm.TaskRun(task)
			return task, nil
		} else {
			return nil, err
		}
	}
}

func (tm *Manager) TaskRun(task *TorrentTask) {
	if err := task.GetInfo(); err != nil {
		log.Printf("获取 Torrent 信息失败 %s \n", err)
	}
}

func (tm *Manager) AddUriTask(uri string) (string, error) {
	// todo 下载参数封装
	t, err := tm.newUriTorrentWithPath(uri, "D:\\Torrent")
	if err != nil {
		return "", err
	}
	_, err = tm.createTask(t)
	if err == nil {
		return t.InfoHash().String(), nil
	}
	return "", err
}

func (tm *Manager) Download(hash string, files []string) error {
	if task := tm.taskMap.Load(hash); task != nil {
		go task.Download(files)
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

func InitTaskManager() {
	tm := GetTaskManager()
	if dbTasks, err := db.SelectActiveTaskList(); err == nil {
		var infoHashList []string
		dbTaskMap := make(map[string]db.Task)
		for _, dbTask := range dbTasks {
			infoHashList = append(infoHashList, dbTask.InfoHash)
			dbTaskMap[dbTask.InfoHash] = dbTask

		}
		if mis, err := db.SelectMateInfoList(infoHashList); err == nil {
			for _, mi := range mis {
				infoHash := mi.HashInfoBytes().String()

				var dbTask db.Task
				if iDbTask, ok := dbTaskMap[infoHash]; ok {
					dbTask = iDbTask
				}
				if err := tm.recoveryTask(&dbTask, mi); err == nil {
					log.Printf("任务 %s 恢复成功 \n", mi.HashInfoBytes().String())
				} else {
					log.Printf("任务 %s 恢复失败 %s \n", mi.HashInfoBytes().String(), err)
				}
			}
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
}
