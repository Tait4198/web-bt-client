package task

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/web-bt-client/db"
	"github.com/web-bt-client/ws"
	"log"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	client    *torrent.Client
	taskMap   *Map
	execQueue *ExecQueue
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
	if t, n, err := tm.client.AddTorrentSpec(spec); err == nil {
		if !n {
			return nil, fmt.Errorf("任务 %s 已在客户端无法创建", t.InfoHash().String())
		}
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

	if t, n, err := tm.client.AddTorrentSpec(spec); err == nil {
		if !n {
			return nil, fmt.Errorf("任务 %s 已在客户端无法创建", t.InfoHash().String())
		}
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
			Complete:          false,
			Pause:             false,
			MetaInfo:          false,
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
			// 设置Torrent信息 磁力链接/种子文件
			param.createTorrentInfo = createTorrentInfo
			task := newTask(t, tm, param)
			tm.taskMap.Store(t.InfoHash().String(), task)

			ws.GetWebSocketManager().Broadcast(TorrentDbTask{
				Type:  Add,
				Wait:  false,
				Queue: false,
				Task:  dbTask,
			})

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

func (tm *Manager) addToQueue(task *TorrentTask, reloadTorrent bool) error {
	infoHash := task.param.InfoHash
	if _, index := tm.execQueue.find(infoHash); index != -1 {
		return fmt.Errorf("任务 %s 已在队列中", infoHash)
	}
	if err := task.ready(reloadTorrent); err != nil {
		return err
	}
	tm.execQueue.pushBack(task)
	BroadcastTaskStatus(task, QueueStatus, true)
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
		t.Drop()
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
	if err := tm.addToQueue(task, false); err != nil {
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
		return tm.addToQueue(task, true)
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
	if err := task.stop(); err != nil {
		return err
	}

	// Queue Remove
	tm.execQueue.remove(task.param.InfoHash)

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
		return task.GetTorrentInfo(true)
	}
	if task, err := tm.recoveryTaskWithHash(infoHash); err == nil {
		return task.GetTorrentInfo(true)
	}
	return TorrentInfoWrapper{}, fmt.Errorf("未找到 %s 信息", infoHash)
}

func (tm *Manager) GetTasks() ([]TorrentDbTask, error) {
	if dbTasks, err := db.SelectTaskList(); err != nil {
		return nil, err
	} else {
		if dbTasks == nil {
			return nil, err
		}
		var tasks []TorrentDbTask
		for _, dbTask := range dbTasks {
			if task, err := tm.getTask(dbTask.InfoHash); err == nil &&
				task.torrent.BytesCompleted() > dbTask.CompleteFileLength {
				dbTask.CompleteFileLength = task.torrent.BytesCompleted()
			}
			tdTask := TorrentDbTask{
				Task: dbTask,
			}
			// 是否等待状态
			if task, err := tm.getTask(dbTask.InfoHash); err == nil {
				tdTask.Wait = task.wait
			}
			// 是否在队列
			if _, index := tm.execQueue.find(dbTask.InfoHash); index > 0 {
				tdTask.Queue = true
			} else {
				tdTask.Queue = false
			}
			tasks = append(tasks, tdTask)
		}
		return tasks, err
	}
}

var taskManager *Manager
var tmOnce sync.Once

func GetTaskManager() *Manager {
	tmOnce.Do(func() {
		cfg := torrent.NewDefaultClientConfig()
		cfg.Seed = true
		//cfg.Logger = logger.Discard
		cfg.ListenPort = 42077
		client, err := torrent.NewClient(cfg)
		if err != nil {
			log.Fatalln(err)
		}
		taskManager = &Manager{
			client:    client,
			taskMap:   NewTaskMap(),
			execQueue: NewExecQueue(1),
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
					if err := tm.addToQueue(task, false); err == nil {
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
