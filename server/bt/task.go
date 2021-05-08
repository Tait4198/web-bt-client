package bt

import (
	"fmt"
	"github.com/anacrolix/torrent"
)

type TorrentTask struct {
	torrent  *torrent.Torrent
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
			t.DownloadAll()
			dt.download.downloadEnd <- struct{}{}
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

		stop := false
		select {
		case <-dt.torrent.GotInfo():
			break
		case <-dt.info.stop:
			stop = true
			break
		}
		if !stop {
			info := dt.torrent.Info()
			fmt.Println(info.Name)
		} else {
			fmt.Printf("Stop Get Info %s\n", dt.torrent.InfoHash().String())
		}
		dt.info.run = false
		// 通知信息已获取
		dt.info.getInfoEnd <- struct{}{}
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

func NewTorrentTask(t *torrent.Torrent) *TorrentTask {
	return &TorrentTask{
		torrent: t,
		info: infoStatus{
			stop:       make(chan struct{}),
			getInfoEnd: make(chan struct{}),
			run:        false,
		},
		download: downloadStatus{
			stop: make(chan struct{}),
			run:  false,
		},
	}
}
