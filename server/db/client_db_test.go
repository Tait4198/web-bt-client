package db

import (
	"context"
	"crawshaw.io/sqlite/sqlitex"
	"fmt"
	"log"
	"testing"
)

func TestDbInit(t *testing.T) {
	InitDb()
	pool := GetPool()
	defer pool.Close()
	conn := pool.Get(context.TODO())
	defer pool.Put(conn)

	err := sqlitex.ExecScript(conn, `
		create table if not exists task_list
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

	stmt := conn.Prep("select count(*) from task_list where info_hash = $hash")
	stmt.SetText("$hash", "4adb90ece042c4d38446cc0a3954d043091ababf")
	if v, err := sqlitex.ResultInt(stmt); err == nil {
		fmt.Println(v)
	} else {
		log.Fatalln(err)
	}
}
