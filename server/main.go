package main

import (
	"github.com/web-bt-client/db"
	btHttp "github.com/web-bt-client/http"
	"net/http"
)

func main() {
	db.InitDb()
	http.ListenAndServe(":8080", btHttp.Router())
}
