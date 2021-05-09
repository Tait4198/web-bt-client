package db

import (
	"bufio"
	"context"
	"crawshaw.io/sqlite/sqlitex"
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"log"
	"os"
	"testing"
)

func TestDbInit(t *testing.T) {
	InitDb()
	pool := GetPool()
	defer pool.Close()
	conn := pool.Get(context.TODO())
	defer pool.Put(conn)

	err := sqlitex.ExecScript(conn, `
		create table if not exists tasks
		(
			id        integer primary key autoincrement,
			info_hash text not null ,
			unique (info_hash)
		);
	`)
	if err != nil {
		log.Println(err)
	}
}

func TestQueryCount(t *testing.T) {
	InitDb()
	pool := GetPool()
	defer pool.Close()
	conn := pool.Get(context.TODO())
	defer pool.Put(conn)

	stmt := conn.Prep("select count(*) from tasks where info_hash = $hash")
	stmt.SetText("$hash", "4adb90ece042c4d38446cc0a3954d043091ababf")
	if v, err := sqlitex.ResultInt(stmt); err == nil {
		fmt.Println(v)
	} else {
		log.Fatalln(err)
	}
}

func TestInsertBlob(t *testing.T) {
	InitDb()
	pool := GetPool()
	defer pool.Close()
	conn := pool.Get(context.TODO())
	defer pool.Put(conn)

	b, err := retrieveROM("D:\\Torrent\\ubuntu-20.04.2-live-server-amd64.iso.torrent")
	if err != nil {
		log.Fatalln(err)
	}
	err = sqlitex.Exec(conn, "insert into torrent_data values (?,?);", nil, "", b)
	if err != nil {
		log.Fatalln(err)
	}
}

func TestGetTorrentMetaInfo(t *testing.T) {
	mi, err := GetMetaInfo("90289fd34dfc1cf8f316a268add8354c85334458")
	if err != nil {
		log.Fatalln(err)
	}
	client, _ := torrent.NewClient(nil)
	to, _ := client.AddTorrent(mi)
	fmt.Println(to.Name())
}

func TestReadBlob(t *testing.T) {
	InitDb()
	pool := GetPool()
	defer pool.Close()
	conn := pool.Get(context.TODO())
	defer pool.Put(conn)
	stmt := conn.Prep("select * from torrent_data")

	client, _ := torrent.NewClient(nil)

	for {
		if hasRow, err := stmt.Step(); err != nil {
			// ... handle error
		} else if !hasRow {
			break
		}
		fmt.Println(stmt.ColumnCount())
		r := stmt.GetReader("torrent_blob")
		//b, err := io.ReadAll(r)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//fmt.Println(len(b))

		mi, err := metainfo.Load(r)
		if err != nil {
			log.Fatalln(err)
		}
		t, err := client.AddTorrent(mi)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(t.Info().Name)
	}
}

func retrieveROM(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size = stats.Size()
	bytes := make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	return bytes, err
}
