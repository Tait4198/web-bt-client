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
	stop chan struct{}
	run  bool
}

type downloadStatus struct {
	stop           chan struct{}
	downloadEnd    chan struct{}
	run            bool
	downloadLength int64
}

type TaskTorrentInfo struct {
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`

	Length         int64 `json:"length"`
	BytesCompleted int64 `json:"bytes_completed"`

	Pieces          int `json:"pieces"`
	CompletedPieces int `json:"completed_pieces"`

	Files []TaskTorrentInfoFile `json:"files"`
}

type TaskTorrentInfoFile struct {
	Path []string `json:"path"`

	Length         int64 `json:"length"`
	BytesCompleted int64 `json:"bytes_completed"`

	Pieces          int `json:"pieces"`
	CompletedPieces int `json:"completed_pieces"`
}

func (dt *TorrentTask) GetTaskTorrentInfo() TaskTorrentInfo {
	t := dt.torrent
	torrentInfo := TaskTorrentInfo{
		InfoHash:       t.InfoHash().String(),
		Name:           t.Name(),
		Length:         t.Length(),
		BytesCompleted: t.BytesCompleted(),
	}

	completedPieces := 0
	for _, psr := range t.PieceStateRuns() {
		if psr.Complete {
			completedPieces += psr.Length
		}
	}
	torrentInfo.CompletedPieces = completedPieces
	torrentInfo.Pieces = t.NumPieces()

	var files = make([]TaskTorrentInfoFile, len(t.Files()))
	for i, f := range t.Files() {
		filePieces := 0
		fileCompletedPieces := 0
		for _, state := range f.State() {
			filePieces++
			if state.Complete {
				fileCompletedPieces++
			}
		}
		torrentFile := TaskTorrentInfoFile{
			Path:            f.FileInfo().Path,
			Length:          f.Length(),
			BytesCompleted:  f.BytesCompleted(),
			Pieces:          filePieces,
			CompletedPieces: fileCompletedPieces,
		}
		files[i] = torrentFile
	}
	if len(files) == 1 && files[0].Path == nil {
		files[0].Path = []string{t.Name()}
	}

	torrentInfo.Files = files

	return torrentInfo
}

func (dt *TorrentTask) PrintStatus() {
	t := dt.torrent

	line2 := fmt.Sprintf("BytesRead %d BytesReadData(下载速度) %d BytesReadUsefulData %d",
		t.Stats().BytesRead,
		t.Stats().BytesReadData,
		t.Stats().BytesReadUsefulData,
	)
	fmt.Println(line2)

	line3 := fmt.Sprintf("BytesWritten %d BytesWrittenData %d",
		t.Stats().BytesWritten,
		t.Stats().BytesWrittenData,
	)
	fmt.Println(line3)

	line4 := fmt.Sprintf("ChunksRead %d ChunksReadUseful %d ChunksReadWasted %d ChunksWritten %d MetadataChunksRead %d",
		t.Stats().ChunksRead,
		t.Stats().ChunksReadUseful,
		t.Stats().ChunksReadWasted,
		t.Stats().ChunksWritten,
		t.Stats().MetadataChunksRead,
	)
	fmt.Println(line4)

	line5 := fmt.Sprintf("TotalPeers %d ActivePeers %d HalfOpenPeers %d PendingPeers %d ",
		t.Stats().TotalPeers,
		t.Stats().ActivePeers,
		t.Stats().HalfOpenPeers,
		t.Stats().PendingPeers,
	)
	fmt.Println(line5)
}

func (dt *TorrentTask) Download(files []string) {
	go func() {
		t := dt.torrent
		if t.Info() == nil {
			log.Printf("下载开始前需要先获取种子信息")
			return
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
		download:
			for {
				select {
				case <-time.After(time.Second):
					downloadEnd := true
					for _, f := range t.Files() {
						fmt.Printf("%d - %d\n", f.BytesCompleted(), f.Length())
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
						dt.PrintStatus()
					} else {
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
		dt.download.run = false
		log.Printf("Download End %s %s \n", t.Name(), t.InfoHash().String())
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
				log.Printf("完成获取 %s \n", infoHash)
			} else {
				log.Printf("停止获取 %s \n", infoHash)
			}
		} else if torrentCount == -1 {
			log.Printf("获取Torrent数量失败 %s \n", infoHash)
		} else {
			if mi, err := db.GetMetaInfo(infoHash); err == nil {
				// drop torrent
				t.Drop()
				if nt, err := dt.client.AddTorrent(mi); err == nil {
					dt.torrent = nt
					log.Printf("完成获取 %s \n", infoHash)
				} else {
					log.Println(fmt.Errorf("MetaInfo 转换 Torrent 失败 %w \n", err))
				}
			} else {
				log.Println(fmt.Errorf("GetMetaInfo 失败 %w \n", err))
			}
		}
		dt.info.run = false
		log.Printf("GetInfo End %s %s \n", t.Name(), t.InfoHash().String())
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
			stop: make(chan struct{}),
			run:  false,
		},
		download: downloadStatus{
			stop:        make(chan struct{}),
			downloadEnd: make(chan struct{}),
			run:         false,
		},
	}
}
