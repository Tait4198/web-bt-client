package task

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/web-bt-client/db"
	"github.com/web-bt-client/ws"
	"log"
	"time"
)

type TorrentTask struct {
	torrent  *torrent.Torrent
	manager  *Manager
	active   bool
	info     infoStatus
	download downloadStatus
	param    Param
	wait     bool
}

type infoStatus struct {
	stop chan struct{}
	run  bool
}

type downloadStatus struct {
	stop                chan struct{}
	downloadEnd         chan struct{}
	run                 bool
	downloadLength      int64
	lastCompletedPieces int
}

func (dt *TorrentTask) GetTaskParam() *Param {
	return &dt.param
}

func (dt *TorrentTask) taskDownload() {
	if dt.download.run || !dt.active {
		return
	}
	t := dt.torrent
	defer func() {
		dt.download.run = false
		log.Printf("Download End %s %s \n", t.Name(), t.InfoHash().String())
	}()

	if t.Info() == nil {
		if err := dt.taskGetInfo(); err != nil {
			log.Printf("获取 Torrent 信息失败 %s \n", err)
			return
		}
	}
	dt.download.run = true
	go func() {
		fMap := make(map[string]byte)
		if !dt.param.DownloadAll {
			for _, file := range dt.param.DownloadFiles {
				fMap[file] = 0
			}
		}
		for _, f := range t.Files() {
			if dt.param.DownloadAll {
				f.Download()
				dt.download.downloadLength += f.Length()
			} else {
				if _, ok := fMap[f.DisplayPath()]; ok {
					f.Download()
					dt.download.downloadLength += f.Length()
				} else {
					f.SetPriority(torrent.PiecePriorityNone)
				}
			}
		}
		log.Printf("Start Download %s %s \n", t.Name(), t.InfoHash().String())
	}()

	go func() {
		wsm := ws.GetWebSocketManager()
	download:
		for {
			select {
			case <-time.After(time.Second):
				downloadEnd := true
				for _, f := range t.Files() {
					if f.Priority() != torrent.PiecePriorityNone && f.BytesCompleted() != f.Length() {
						downloadEnd = false
					}
				}

				if !downloadEnd {
					var completedPieces, partialPieces int
					psrs := t.PieceStateRuns()
					for _, r := range psrs {
						if r.Complete {
							completedPieces += r.Length
						}
						if r.Partial {
							partialPieces += r.Length
						}
					}

					if completedPieces > dt.download.lastCompletedPieces {
						if err := db.UpdateTaskCompleteFileLength(t.BytesCompleted(), t.InfoHash().String()); err == nil {
							dt.download.lastCompletedPieces = completedPieces
						} else {
							log.Printf("任务 %s CompleteFileLength 更新失败 %s", t.InfoHash().String(), err)
						}
					}

					wsm.Broadcast(dt.GetTorrentStats(false, true, true))
				} else {
					if err := db.TaskComplete(t.BytesCompleted(), t.InfoHash().String()); err != nil {
						log.Printf("任务 %s 完成信息更新失败 %s \n", t.InfoHash().String(), err)
					}
					wsm.Broadcast(TorrentTaskComplete{
						TorrentTaskStatus: TorrentTaskStatus{
							TorrentBase: TorrentBase{
								InfoHash: dt.param.InfoHash,
								Type:     Complete,
							},
							Status: true,
						},
						LastCompleteLength: dt.torrent.BytesCompleted(),
					})
					break download
				}
			case <-dt.download.stop:
				log.Printf("Download Stop %s %s \n", t.Name(), t.InfoHash().String())
				break download
			}
		}
		dt.download.downloadEnd <- struct{}{}
	}()

	<-dt.download.downloadEnd
}

func (dt *TorrentTask) taskGetInfo() error {
	t := dt.torrent

	defer func() {
		if t.Info() != nil {
			fileLen := t.Length()
			if fileLen == 0 {
				for _, f := range t.Files() {
					fileLen += f.Length()
				}
			}
			if err := db.UpdateTaskMetaInfo(t.InfoHash().String(), t.Info().Name, fileLen); err != nil {
				log.Printf("任务 %s 更新 MetaInfo 失败 \n", t.InfoHash().String())
			}
			if dt.param.DownloadAll {
				var downloadFiles []string
				for _, f := range t.Files() {
					downloadFiles = append(downloadFiles, f.DisplayPath())
				}
				if err := db.UpdateTaskDownloadFiles(downloadFiles, t.InfoHash().String()); err != nil {
					log.Printf("任务 %s 更新 DownloadFiles 失败 \n", t.InfoHash().String())
				}
			}
		}
		dt.info.run = false
		log.Printf("GetInfo End %s %s \n", t.Name(), t.InfoHash().String())
	}()

	if dt.info.run || dt.torrent.Info() != nil || !dt.active {
		return nil
	}
	infoHash := dt.torrent.InfoHash().String()
	dt.info.run = true

	torrentCount := db.SelectTorrentDataCount(infoHash)
	if torrentCount == 0 {
		log.Printf("开始获取 %s \n", infoHash)
		stop := false
		select {
		case <-dt.torrent.GotInfo():
			break
		case <-dt.info.stop:
			stop = true
			break
		}
		if !stop {
			if b, err := bencode.Marshal(t.Metainfo()); err == nil {
				if err := db.InsertTorrentData(infoHash, b); err != nil {
					return fmt.Errorf("Hash %s 写入 SQLite 失败 %w \n", infoHash, err)
				}
				if ti, err := dt.GetTorrentInfo(false); err == nil {
					ws.GetWebSocketManager().Broadcast(ti)
				} else {
					return fmt.Errorf(" %s 信息完成获取推送失败 %w", infoHash, err)
				}
			} else {
				return fmt.Errorf("bencode.Marshal 失败 %w \n", err)
			}
		} else {
			return fmt.Errorf("主动停止获取 %s \n", infoHash)
		}
	} else if torrentCount == -1 {
		return fmt.Errorf("获取Torrent数量失败 %s \n", infoHash)
	} else {
		t.Drop()
		if nt, err := dt.manager.recoveryTorrentWithHash(infoHash); err == nil {
			dt.torrent = nt
			t = nt
		} else {
			return err
		}
	}

	// GetInfo Success
	log.Printf("完成获取 %s \n", infoHash)
	return nil
}

func (dt *TorrentTask) stop() error {
	// 停止Task
	if !dt.active {
		return fmt.Errorf("任务 %s 已停止", dt.param.InfoHash)
	}
	if err := db.UpdateTaskPause(true, dt.param.InfoHash); err != nil {
		return fmt.Errorf("任务 %s 停止失败 %w", dt.torrent.InfoHash().String(), err)
	}
	if dt.info.run {
		dt.info.stop <- struct{}{}
	}
	if dt.download.run {
		dt.download.stop <- struct{}{}
	}
	dt.active = false
	dt.torrent.Drop()
	return nil
}

func (dt *TorrentTask) ready(reloadTorrent bool) error {
	if dt.active {
		return fmt.Errorf("任务 %s 已启动", dt.param.InfoHash)
	}
	if reloadTorrent {
		if err := dt.reloadTorrent(); err != nil {
			return err
		}
	}
	if err := db.UpdateTaskPause(false, dt.param.InfoHash); err != nil {
		return fmt.Errorf("任务 %s 启动失败 %w", dt.torrent.InfoHash().String(), err)
	}
	dt.active = true
	return nil
}

func (dt *TorrentTask) exec() {
	if err := dt.taskGetInfo(); err == nil && dt.param.Download {
		dt.taskDownload()
	}
	dt.active = false
	if err := db.UpdateTaskPause(true, dt.torrent.InfoHash().String()); err == nil {
		BroadcastTaskStatus(dt, Pause, true)
	}
}

func (dt *TorrentTask) reloadTorrent() error {
	dt.torrent.Drop()
	var mi *metainfo.MetaInfo
	infoHash := dt.torrent.InfoHash().String()
	if dt.torrent.Info() == nil {
		// 种子文件任务必定存在MetaInfo
		if m, err := db.SelectMetaInfo(infoHash); err == nil {
			mi = m
		} else {
			// 磁力链接任务
			if nt, err := dt.manager.newUriTorrentWithPath(dt.param.CreateTorrentInfo, dt.param.DownloadPath); err == nil {
				dt.torrent = nt
				return nil
			} else {
				return fmt.Errorf("任务 %s 重新加载 Torrent 失败 (磁力链接)%w", infoHash, err)
			}
		}
	} else {
		m := dt.torrent.Metainfo()
		mi = &m
	}
	if nt, err := dt.manager.newMetaInfoTorrentWithPath(mi, dt.param.DownloadPath); err == nil {
		dt.torrent = nt
	} else {
		return fmt.Errorf("任务 %s 重新加载 Torrent 失败 (MetaInfo)%w", infoHash, err)
	}
	return nil
}

func (dt *TorrentTask) TaskWait() error {
	if !dt.wait {
		dt.wait = true
		BroadcastTaskStatus(dt, Wait, dt.wait)

		go func() {
			tick := time.Tick(time.Second * 3)
			<-tick

			dt.wait = false
			BroadcastTaskStatus(dt, Wait, dt.wait)
		}()
		return nil
	}
	return fmt.Errorf("任务 %s 等待中", dt.torrent.InfoHash().String())
}

func (dt *TorrentTask) FileComplete(path string) bool {
	for _, f := range dt.torrent.Files() {
		fmt.Println(f.DisplayPath())
		if f.DisplayPath() == path {
			return true
		}
	}
	return false
}

func (dt *TorrentTask) GetTorrentName() string {
	return dt.torrent.Name()
}

func newTask(t *torrent.Torrent, tm *Manager, param Param) *TorrentTask {
	return &TorrentTask{
		torrent: t,
		manager: tm,
		active:  false,
		info: infoStatus{
			stop: make(chan struct{}),
			run:  false,
		},
		download: downloadStatus{
			stop:        make(chan struct{}),
			downloadEnd: make(chan struct{}),
			run:         false,
		},
		param: param,
		wait:  false,
	}
}
