package http

import (
	"github.com/gin-gonic/gin"
	"github.com/web-bt-client/task"
	"net/http"
)

func torrentInfo(c *gin.Context) {
	hash := c.DefaultQuery("hash", "")
	torrentInfo, err := task.GetTaskManager().GetTorrentInfo(hash)
	if err == nil {
		c.JSON(http.StatusOK, DataJson(true, torrentInfo))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func InitTorrentRouter(groupRouter *gin.RouterGroup) {
	groupRouter.GET("/info", torrentInfo)
}
