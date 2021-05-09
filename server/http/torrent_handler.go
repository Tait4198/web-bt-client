package http

import (
	"github.com/gin-gonic/gin"
	"github.com/web-bt-client/bt"
	"net/http"
)

type DownloadParam struct {
	InfoHash string   `json:"hash"`
	Files    []string `json:"files"`
}

func downloadTorrent(c *gin.Context) {
	dp := DownloadParam{}
	err := c.BindJSON(&dp)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageJson(false, "无效参数"))
		return
	}
	if err := bt.GetTaskManager().Download(dp.InfoHash, dp.Files); err == nil {
		c.JSON(http.StatusOK, DataJson(true, dp.InfoHash))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func InitTorrentRouter(groupRouter *gin.RouterGroup) {
	groupRouter.POST("/download", downloadTorrent)
}
