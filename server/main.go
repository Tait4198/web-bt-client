package main

import (
	btHttp "github.com/web-bt-client/http"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", btHttp.Router())
}
