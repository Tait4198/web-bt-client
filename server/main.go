package main

import (
	"github.com/web-bt-client/db"
	btHttp "github.com/web-bt-client/http"
	"github.com/web-bt-client/task"
	"net/http"
)

func main() {
	db.InitDb()
	task.InitTaskManager()
	http.ListenAndServe(":8080", btHttp.Router())
}
