package bt

import (
	"fmt"
	"github.com/anacrolix/torrent"
)

type TorrentTask struct {
	torrent     *torrent.Torrent
	stopGetInfo chan struct{}
}

func (dt *TorrentTask) GetInfo() {
	go func() {
		stop := false
		select {
		case <-dt.torrent.GotInfo():
			break
		case <-dt.stopGetInfo:
			stop = true
			break
		}
		if !stop {
			info := dt.torrent.Info()
			fmt.Println(info.Name)
		} else {
			fmt.Printf("Stop Get Info %s\n", dt.torrent.InfoHash().String())
		}
		dt.torrent.Drop()
	}()
}

func (dt *TorrentTask) StopGetInfo() {
	dt.stopGetInfo <- struct{}{}
}

func NewTorrentTask(t *torrent.Torrent) *TorrentTask {
	return &TorrentTask{
		torrent:     t,
		stopGetInfo: make(chan struct{}),
	}
}
