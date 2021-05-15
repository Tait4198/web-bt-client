package task

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/web-bt-client/db"
	"github.com/web-bt-client/ws"
	"log"
	"time"
)

type TorrentTask struct {
	torrent      *torrent.Torrent
	manager      *Manager
	info         infoStatus
	downloadPath string
	download     downloadStatus
}

type infoStatus struct {
	stop chan struct{}
	run  bool
}

type downloadStatus struct {
	stop           chan struct{}
	downloadEnd    chan struct{}
	run            bool
	downloadLength int64
}

func (dt *TorrentTask) Download(files []string) {
	t := dt.torrent
	defer func() {
		dt.download.run = false
		log.Printf("Download End %s %s \n", t.Name(), t.InfoHash().String())
	}()

	if t.Info() == nil {
		if err := dt.GetInfo(); err != nil {
			log.Printf("获取 Torrent 信息失败 %s \n", err)
			return
		}
	}

	dt.download.run = true
	go func() {
		fMap := make(map[string]byte)
		for _, file := range files {
			fMap[file] = 0
		}
		for _, f := range t.Files() {
			if len(files) == 0 {
				f.Download()
			} else {
				if _, ok := fMap[f.DisplayPath()]; ok {
					f.Download()
				} else {
					f.SetPriority(torrent.PiecePriorityNone)
				}
			}
			dt.download.downloadLength += f.Length()
		}
		log.Printf("Start Download %s %s \n", t.Name(), t.InfoHash().String())
	}()

	go func() {
		start := time.Now()
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
					line := fmt.Sprintf(
						"\n%v: downloading %q: %d/%d, %d/%d pieces completed (%d partial)",
						time.Since(start),
						t.Name(),
						uint64(t.BytesCompleted()),
						uint64(t.Length()),
						completedPieces,
						t.NumPieces(),
						partialPieces,
					)
					fmt.Println(line)

					wsm.Broadcast(dt.GetTorrentStats())
				} else {
					// Download End
					t.Drop()
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

func (dt *TorrentTask) GetInfo() error {
	if dt.torrent.Info() != nil {
		return nil
	}

	dt.info.run = true
	infoHash := dt.torrent.InfoHash().String()
	t := dt.torrent

	defer func() {
		dt.info.run = false
		log.Printf("GetInfo End %s %s \n", t.Name(), t.InfoHash().String())
	}()

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
			} else {
				return fmt.Errorf("bencode.Marshal 失败 %w \n", err)
			}
		} else {
			return fmt.Errorf("主动停止获取 %s \n", infoHash)
		}
	} else if torrentCount == -1 {
		return fmt.Errorf("获取Torrent数量失败 %s \n", infoHash)
	} else {
		if mi, err := db.SelectMetaInfo(infoHash); err == nil {
			// drop torrent
			t.Drop()
			if nt, err := dt.manager.newMetaInfoTorrentWithPath(mi, dt.downloadPath); err == nil {
				dt.torrent = nt
				// update t
				t = nt
			} else {
				return fmt.Errorf("MetaInfo 转换 Torrent 失败 %w \n", err)
			}
		} else {
			return fmt.Errorf("GetMetaInfo 失败 %w \n", err)
		}
	}

	if err := db.UpdateTaskMetaInfo(t.InfoHash().String(), t.Info().Name, t.Info().Length); err != nil {
		return err
	}

	// GetInfo Success
	log.Printf("完成获取 %s \n", infoHash)
	return nil
}

func (dt *TorrentTask) Stop() {
	if dt.info.run {
		dt.info.stop <- struct{}{}
	}
	if dt.download.run {
		dt.download.stop <- struct{}{}
	}
	dt.torrent.Drop()
}

func (dt *TorrentTask) Start() error {
	mi := dt.torrent.Metainfo()
	if nt, err := dt.manager.client.AddTorrent(&mi); err == nil {
		dt.torrent = nt

		dt.Download([]string{})

		return nil
	} else {
		return fmt.Errorf("任务开始失败 %w", err)
	}
}

func newTask(t *torrent.Torrent, tm *Manager) *TorrentTask {
	return newTaskSetPath(t, tm, "")
}

func newTaskSetPath(t *torrent.Torrent, tm *Manager, downloadPath string) *TorrentTask {
	return &TorrentTask{
		torrent: t,
		manager: tm,
		info: infoStatus{
			stop: make(chan struct{}),
			run:  false,
		},
		download: downloadStatus{
			stop:        make(chan struct{}),
			downloadEnd: make(chan struct{}),
			run:         false,
		},
		downloadPath: downloadPath,
	}
}
