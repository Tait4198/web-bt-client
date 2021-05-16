package db

import (
	"context"
	"crawshaw.io/sqlite/sqlitex"
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
			info_hash            char(40) primary key not null,
			torrent_name         varchar(2048)        not null,
			complete             tinyint(1)                    default 0,
			meta_info            tinyint(1)                    default 0,
			pause                tinyint(1)                    default 0,
			download             tinyint(1)                    default 0,
			download_path        varchar(2048)        not null,
			download_files       text                 not null default '',
			file_length          bigint                        default 0,
			complete_file_length bigint                        default 0,
			create_time          datetime             not null default current_timestamp,
			complete_time        datetime,
			create_torrent_info  varchar(2048)        not null
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

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
