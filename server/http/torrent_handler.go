package http

import (
	"fmt"
	"github.com/anacrolix/torrent/metainfo"
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

func torrentUpload(c *gin.Context) {
	file, err := c.FormFile("torrent")
	if err != nil {
		c.JSON(http.StatusOK, MessageJson(false, fmt.Sprintf("文件上传失败 %s", err)))
		return
	}
	fo, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, MessageJson(false, fmt.Sprintf("文件打开失败 %s", err)))
		return
	}
	defer fo.Close()
	mi, err := metainfo.Load(fo)
	if err != nil {
		c.JSON(http.StatusOK, MessageJson(false, fmt.Sprintf("metainfo.Load 失败 %s", err)))
		return
	}
	if tiw, err := task.GetTaskManager().SaveTorrent(mi); err == nil {
		c.JSON(http.StatusOK, DataJson(true, tiw))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func InitTorrentRouter(groupRouter *gin.RouterGroup) {
	groupRouter.GET("/info", torrentInfo)
	groupRouter.POST("/upload", torrentUpload)
}
