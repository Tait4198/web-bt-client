package db

import (
	"context"
	"crawshaw.io/sqlite/sqlitex"
	"fmt"
	"github.com/anacrolix/torrent/metainfo"
	"log"
	"sync"
)

var pool *sqlitex.Pool
var poolOnce sync.Once

func InitDb() {
	conn := GetPool().Get(context.TODO())
	defer GetPool().Put(conn)
	err := sqlitex.ExecScript(conn, `
		create table if not exists torrent_data
		(
			info_hash    char(40) not null,
			torrent_blob blob,
			primary key (info_hash)
		);

		create table if not exists tasks
		(
			id        integer primary key autoincrement,
			info_hash text not null ,
			unique (info_hash)
		);
	`)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("DB INIT")
}

func GetPool() *sqlitex.Pool {
	poolOnce.Do(func() {
		inPool, err := sqlitex.Open("file:.web-bt-client.db", 0, 10)
		if err != nil {
			log.Fatal(err)
		}
		pool = inPool
	})
	return pool
}

func ExecSql(sql string, args ...interface{}) error {
	conn := GetPool().Get(context.TODO())
	defer GetPool().Put(conn)
	return sqlitex.Exec(conn, sql, nil, args...)
}

func GetMetaInfo(infoHash string) (*metainfo.MetaInfo, error) {
	conn := GetPool().Get(context.TODO())
	defer GetPool().Put(conn)
	stmt := conn.Prep("select * from torrent_data where info_hash = $hash")
	defer func() {
		if err := stmt.Reset(); err != nil {
			log.Println(fmt.Errorf("GetMetaInfo stmt.Reset() 失败 %w \n", err))
		}
	}()
	stmt.SetText("$hash", infoHash)
	if hasRow, err := stmt.Step(); err != nil {
		return nil, err
	} else if !hasRow {
		return nil, fmt.Errorf(" Hash %s 不存在 ", infoHash)
	}
	r := stmt.GetReader("torrent_blob")
	if mi, err := metainfo.Load(r); err == nil {
		return mi, nil
	} else {
		return nil, fmt.Errorf("bytes.Reader 转换 MetaInfo 失败 %w \n", err)
	}
}

func GetTorrentDataCount(infoHash string) int {
	conn := GetPool().Get(context.TODO())
	defer GetPool().Put(conn)
	stmt := conn.Prep("select count(*) from torrent_data where info_hash = $hash")
	stmt.SetText("$hash", infoHash)
	if count, err := sqlitex.ResultInt(stmt); err == nil {
		return count
	} else {
		log.Println(fmt.Errorf("获取 Hash %s 数量失败 %w \n", infoHash, err))
	}
	return -1
}
