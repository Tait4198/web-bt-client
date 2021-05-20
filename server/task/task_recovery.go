package task

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/web-bt-client/db"
	"regexp"
)

func (tm *Manager) recoveryTorrent(dbTask *db.Task, mi *metainfo.MetaInfo) (*torrent.Torrent, error) {
	var t *torrent.Torrent
	if mi != nil {
		// MetaInfo恢复
		if mt, err := tm.newMetaInfoTorrentWithPath(mi, dbTask.DownloadPath); err == nil {
			t = mt
		} else {
			return nil, err
		}
	} else if match, _ := regexp.MatchString("magnet:\\?xt=urn:btih:[a-z0-9]{40}.*", dbTask.CreateTorrentInfo); match {
		// 磁力链接恢复
		if mt, err := tm.newUriTorrentWithPath(dbTask.CreateTorrentInfo, dbTask.DownloadPath); err == nil {
			t = mt
		} else {
			return nil, err
		}
	} else {
		// todo 文件恢复
		return nil, fmt.Errorf("文件恢复未实现")
	}
	return t, nil
}

func (tm *Manager) recoveryTorrentWithHash(infoHash string) (*torrent.Torrent, error) {
	// 从数据库中恢复
	dbTask, err := db.SelectTask(infoHash)
	if err != nil {
		return nil, fmt.Errorf("任务 %s 信息不存在", infoHash)
	}
	mi, _ := db.SelectMetaInfo(infoHash)
	return tm.recoveryTorrent(&dbTask, mi)
}

func (tm *Manager) recoveryTask(dbTask *db.Task, mi *metainfo.MetaInfo) (*TorrentTask, error) {
	t, err := tm.recoveryTorrent(dbTask, mi)
	if err != nil {
		return nil, err
	}
	if t != nil {
		infoHash := t.InfoHash().String()
		// 恢复下载进度
		if t.BytesCompleted() > 0 && t.BytesCompleted() != dbTask.CompleteFileLength {
			if err := db.UpdateTaskCompleteFileLength(t.BytesCompleted(), infoHash); err != nil {
				return nil, fmt.Errorf("任务 %s 下载进度恢复失败", infoHash)
			}
		}
		param := Param{
			InfoHash:      infoHash,
			Download:      dbTask.Download,
			DownloadFiles: dbTask.DownloadFiles,
			DownloadPath:  dbTask.DownloadPath,
		}
		task := newTask(t, tm, param)
		tm.taskMap.Store(t.InfoHash().String(), task)
		return task, nil
	} else {
		return nil, fmt.Errorf("任务 %s 恢复失败", dbTask.InfoHash)
	}
}

func (tm *Manager) recoveryTaskWithHash(infoHash string) (*TorrentTask, error) {
	// 从数据库中恢复
	dbTask, err := db.SelectTask(infoHash)
	if err != nil {
		return nil, fmt.Errorf("任务 %s 信息不存在", infoHash)
	}
	mi, _ := db.SelectMetaInfo(infoHash)
	return tm.recoveryTask(&dbTask, mi)
}
