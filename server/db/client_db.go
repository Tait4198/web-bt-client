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
		create table if not exists task_list
		(
			id        integer primary key autoincrement,
			info_hash char(40) not null ,
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
