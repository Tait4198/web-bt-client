package task

import (
	"fmt"
	logger "github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/web-bt-client/db"
	"log"
	"strings"
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

func (tm *Manager) createTask(t *torrent.Torrent, createTorrentInfo string, param Param) (*TorrentTask, error) {
	if tm.taskExists(t.InfoHash().String()) {
		return nil, fmt.Errorf("任务已存在 %s", t.InfoHash().String())
	} else {
		dbTask := db.Task{
			InfoHash:          t.InfoHash().String(),
			DownloadPath:      param.DownloadPath,
			Download:          param.Download,
			DownloadFiles:     param.DownloadFiles,
			CreateTime:        time.Now().UnixNano(),
			CreateTorrentInfo: createTorrentInfo,
		}
		if t.Info() == nil {
			dbTask.TorrentName = t.InfoHash().String()
		} else {
			dbTask.TorrentName = t.Info().Name
		}
		if err := db.InsertTask(dbTask); err == nil {
			task := newTask(t, tm, param)
			tm.taskMap.Store(t.InfoHash().String(), task)
			return task, nil
		} else {
			return nil, err
		}
	}
}

func (tm *Manager) getTask(infoHash string) (*TorrentTask, error) {
	task := tm.taskMap.Load(infoHash)
	if task == nil {
		// 从数据库中恢复
		if tTask, err := tm.recoveryTaskWithHash(infoHash); err == nil {
			task = tTask
		} else {
			return nil, err
		}
	}
	return task, nil
}

func (tm *Manager) TaskRun(task *TorrentTask, reloadTorrent bool) error {
	if err := task.Start(reloadTorrent); err != nil {
		return err
	}
	return nil
}

func (tm *Manager) AddUriTask(uri string, param Param) (string, error) {
	t, err := tm.client.AddMagnet(uri)
	if err != nil {
		return "", fmt.Errorf("无效磁力链接")
	}
	infoHash := t.InfoHash().String()
	if tm.taskMap.Load(infoHash) != nil || db.SelectTaskCount(infoHash) > 0 {
		return "", fmt.Errorf("任务 %s 已存在", t.InfoHash().String())
	}
	if param.DownloadPath != "" {
		t, err = tm.newUriTorrentWithPath(uri, param.DownloadPath)
	} else {
		t, err = tm.client.AddMagnet(uri)
	}
	if err != nil {
		return "", err
	}
	// 设置InfoHash
	param.InfoHash = infoHash
	task, err := tm.createTask(t, strings.ToLower(uri), param)
	if err != nil {
		return "", err
	}
	if err := tm.TaskRun(task, false); err != nil {
		return "", err
	}
	return infoHash, nil
}

func (tm *Manager) Start(param Param, wait bool) error {
	task, err := tm.getTask(param.InfoHash)
	if err != nil {
		return err
	}
	if task != nil {
		if wait {
			if err := task.TaskWait(); err != nil {
				return err
			}
		}
		if param.Update {
			if task.GetTaskParam().Download != param.Download {
				task.GetTaskParam().Download = param.Download
				if err := db.UpdateTaskDownload(task.GetTaskParam().Download, task.param.InfoHash); err != nil {
					return fmt.Errorf("任务 %s 下载状态更新失败 %w", task.param.InfoHash, err)
				}
			}

			if param.DownloadFiles != nil {
				task.GetTaskParam().DownloadFiles = param.DownloadFiles
				if err := db.UpdateTaskDownloadFiles(param.DownloadFiles, task.param.InfoHash); err != nil {
					return fmt.Errorf("任务 %s 下载文件更新失败 %w", task.param.InfoHash, err)
				}
			}
		}
		return tm.TaskRun(task, true)
	}
	return fmt.Errorf("任务 %s 不存在", param.InfoHash)
}

func (tm *Manager) Stop(infoHash string, wait bool) error {
	task := tm.taskMap.Load(infoHash)
	if task == nil {
		return fmt.Errorf("任务 %s 不存在", infoHash)
	}
	if wait {
		if err := task.TaskWait(); err != nil {
			return err
		}
	}
	if err := task.Stop(); err != nil {
		return err
	}
	return nil
}

func (tm *Manager) Restart(param Param) error {
	task, err := tm.getTask(param.InfoHash)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("任务 %s 不存在", param.InfoHash)
	}
	if err := task.TaskWait(); err != nil {
		return err
	}
	if task.active {
		if err := tm.Stop(param.InfoHash, false); err != nil {
			return err
		}
	}
	if err := tm.Start(param, false); err != nil {
		return err
	}
	return nil
}

func (tm *Manager) GetTorrentInfo(infoHash string) (TorrentInfoWrapper, error) {
	if task := tm.taskMap.Load(infoHash); task != nil {
		return task.GetTorrentInfo()
	}
	if task, err := tm.recoveryTaskWithHash(infoHash); err == nil {
		return task.GetTorrentInfo()
	}
	return TorrentInfoWrapper{}, fmt.Errorf("未找到 %s 信息", infoHash)
}

var taskManager *Manager
var tmOnce sync.Once

func GetTaskManager() *Manager {
	tmOnce.Do(func() {
		cfg := torrent.NewDefaultClientConfig()
		cfg.Seed = true
		cfg.Logger = logger.Discard
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
		for _, dbTask := range dbTasks {
			infoHashList = append(infoHashList, dbTask.InfoHash)
		}
		if mis, err := db.SelectMateInfoList(infoHashList); err == nil {
			miMap := make(map[string]*metainfo.MetaInfo)
			for _, mi := range mis {
				infoHash := mi.HashInfoBytes().String()
				miMap[infoHash] = mi
			}
			for _, dbTask := range dbTasks {
				mi := miMap[dbTask.InfoHash]
				if task, err := tm.recoveryTask(&dbTask, mi); err == nil {
					if err := tm.TaskRun(task, false); err == nil {
						log.Printf("任务 %s 恢复成功 \n", dbTask.InfoHash)
					} else {
						log.Println(err)
					}
				} else {
					log.Println(err)
				}
			}
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
	log.Println("TASK INIT")
}
