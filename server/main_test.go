package main

import (
	"container/list"
	"fmt"
	"github.com/anacrolix/torrent"
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
	to, _ := client.AddTorrentFromFile("/Users/chenyitao/Downloads/ubuntu-20.04.2-live-server-amd64.iso.torrent")
	<-to.GotInfo()
	fmt.Println(to.InfoHash().String())
	to.Drop()
}

func TestTick(t *testing.T) {
	timeTickerChan := time.Tick(time.Second * 2)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	<-timeTickerChan
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}

func TestList(t *testing.T) {
	//q := task.NewExecQueue(1)
	//q.PushBack("Ni")
	//q.PushBack("Hao")
	//q.PushBack("A")
	//
	//forever := make(chan bool)
	//<-forever

	l := list.New()
	fmt.Println(l.Front())
}
