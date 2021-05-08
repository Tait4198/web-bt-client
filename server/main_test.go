package main

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"testing"
)

func TestHash(t *testing.T) {
	client, _ := torrent.NewClient(nil)
	torrent, _ := client.AddMagnet("magnet:?xt=urn:btih:4ADB90ECE042C4D38446CC0A3954D043091ABABF")
	fmt.Println(torrent.InfoHash().String())
}
