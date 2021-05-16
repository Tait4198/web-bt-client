package main

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	cfg := torrent.NewDefaultClientConfig()
	fmt.Println(cfg.DataDir)
	client, _ := torrent.NewClient(nil)
	to, _ := client.AddMagnet("magnet:?xt=urn:btih:4ADB90ECE042C4D38446CC0A3954D043091ABABF")
	<-to.GotInfo()
	fmt.Println(to.InfoHash().String())

	data := to.Stats().BytesReadData
	fmt.Println(data.Int64())
}

func TestDownload(t *testing.T) {
	client, _ := torrent.NewClient(nil)
	to, _ := client.AddMagnet("magnet:?xt=urn:btih:4ADB90ECE042C4D38446CC0A3954D043091ABABF")
	<-to.GotInfo()
	to.Drop()
	client.AddTorrentInfoHashWithStorage(to.InfoHash(), storage.NewMMap("D:\\Torrent"))
	fmt.Println(to.Name())
	to.DownloadAll()
	client.WaitAll()
}

func TestTick(t *testing.T) {
	timeTickerChan := time.Tick(time.Second * 2)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	<-timeTickerChan
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
