package main

import (
	"flag"
	"fmt"
	"github.com/web-bt-client/db"
	btHttp "github.com/web-bt-client/http"
	"github.com/web-bt-client/task"
	"net/http"
)

func main() {
	var port, size int

	flag.IntVar(&port, "p", 8080, "服务启动端口")
	flag.IntVar(&size, "q", 5, "最大并行任务限制")
	flag.Parse()

	db.InitDb()
	task.InitTaskManager(size)
	http.ListenAndServe(fmt.Sprintf(":%d", port), btHttp.Router())
}
