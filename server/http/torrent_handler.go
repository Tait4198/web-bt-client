package http

import (
	"github.com/gin-gonic/gin"
	"github.com/web-bt-client/task"
	"net/http"
)

type DownloadParam struct {
	InfoHash string   `json:"hash"`
	Files    []string `json:"files"`
}

func torrentDownload(c *gin.Context) {
	dp := DownloadParam{}
	err := c.BindJSON(&dp)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageJson(false, "无效参数"))
		return
	}
	if err := task.GetTaskManager().Download(dp.InfoHash, dp.Files); err == nil {
		c.JSON(http.StatusOK, DataJson(true, dp.InfoHash))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

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
	groupRouter.POST("/download", torrentDownload)
	groupRouter.GET("/info", torrentInfo)
}
