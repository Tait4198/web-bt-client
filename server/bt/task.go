package bt

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/web-bt-client/db"
	"log"
	"time"
)

type TorrentTask struct {
	torrent  *torrent.Torrent
	client   *torrent.Client
	info     infoStatus
	download downloadStatus
}

type infoStatus struct {
	stop       chan struct{}
	getInfoEnd chan struct{}
	run        bool
}

type downloadStatus struct {
	stop        chan struct{}
	downloadEnd chan struct{}
	run         bool
}

func (dt *TorrentTask) Download() {
	go func() {
		dt.download.run = true

		t := dt.torrent
		if t.Info() == nil {
			dt.GetInfo()
			<-dt.info.getInfoEnd
		}

		go func() {
			log.Printf("开始下载 %s \n", t.InfoHash().String())
			t.DownloadAll()
			dt.download.downloadEnd <- struct{}{}
		}()

		go func() {
			start := time.Now()
			for range time.Tick(time.Second) {
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
					"%v: downloading %q: %d/%d, %d/%d pieces completed (%d partial)\n",
					time.Since(start),
					t.Name(),
					uint64(t.BytesCompleted()),
					uint64(t.Length()),
					completedPieces,
					t.NumPieces(),
					partialPieces,
				)

				fmt.Println(line)
			}
		}()

		select {
		case <-dt.download.downloadEnd:
			break
		case <-dt.info.stop:
			break
		}

		dt.download.run = false
	}()
}

func (dt *TorrentTask) GetInfo() {
	go func() {
		dt.info.run = true
		infoHash := dt.torrent.InfoHash().String()
		t := dt.torrent

		torrentCount := db.GetTorrentDataCount(infoHash)
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
					if err := db.ExecSql("insert into torrent_data values (?,?);", infoHash, b); err != nil {
						log.Println(fmt.Errorf("Hash %s 写入 SQLite 失败 %w \n", infoHash, err))
					}
				} else {
					log.Println(fmt.Errorf("bencode.Marshal 失败 %w \n", err))
				}

				// 通知信息已获取
				dt.info.getInfoEnd <- struct{}{}
				log.Printf("完成获取 %s \n", infoHash)
			} else {
				log.Printf("停止获取 %s \n", infoHash)
			}
		} else if torrentCount == -1 {
			log.Printf("获取Torrent数量失败 %s \n", infoHash)
		} else {
			if mi, err := db.GetMetaInfo(infoHash); err == nil {
				if nt, err := dt.client.AddTorrent(mi); err == nil {
					t.Drop()
					dt.torrent = nt
					// 通知信息已获取
					dt.info.getInfoEnd <- struct{}{}
					log.Printf("完成获取 %s \n", infoHash)
				} else {
					log.Println(fmt.Errorf("MetaInfo 转换 Torrent 失败 %w \n", err))
				}
			} else {
				log.Println(fmt.Errorf("GetMetaInfo 失败 %w \n", err))
			}
		}

		dt.info.run = false
	}()
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

func NewTorrentTask(t *torrent.Torrent, c *torrent.Client) *TorrentTask {
	return &TorrentTask{
		torrent: t,
		client:  c,
		info: infoStatus{
			stop:       make(chan struct{}),
			getInfoEnd: make(chan struct{}, 1),
			run:        false,
		},
		download: downloadStatus{
			stop: make(chan struct{}),
			run:  false,
		},
	}
}
