package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/web-bt-client/bt"
	"net/http"
)

func newTask(c *gin.Context) {
	uri := c.DefaultQuery("uri", "")
	tm := bt.GetTaskManager()
	hash, err := tm.AddUriTask(uri)
	if err == nil {
		fmt.Println(hash)
	}
}

func stopTask(c *gin.Context) {
	hash := c.DefaultQuery("hash", "")
	tm := bt.GetTaskManager()
	tm.StopTask(hash)
}

func Router() http.Handler {
	router := gin.New()
	router.GET("/task/new", newTask)
	router.GET("/task/stop", stopTask)
	return router
}
