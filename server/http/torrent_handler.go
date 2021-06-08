package http

import (
	"fmt"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/gin-gonic/gin"
	"github.com/web-bt-client/base"
	"github.com/web-bt-client/db"
	"github.com/web-bt-client/task"
	"net/http"
	"strings"
)

func torrentInfo(c *gin.Context) {
	hash := c.DefaultQuery("hash", "")
	info, err := task.GetTaskManager().GetTorrentInfo(strings.ToLower(hash))
	if err == nil {
		c.JSON(http.StatusOK, DataJson(true, info))
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

func torrentFileDownload(c *gin.Context) {
	hash := c.DefaultQuery("hash", "")
	path := c.DefaultQuery("path", "")
	t, err := task.GetTaskManager().GetTask(hash)
	if err != nil {
		c.JSON(http.StatusOK, DataJson(false, err.Error()))
		return
	}
	downloadPath := fmt.Sprintf("%s/%s", t.GetTaskParam().DownloadPath, t.GetTorrentName())
	if base.Exists(downloadPath) {
		var fullPath string
		if base.IsDir(downloadPath) {
			fullPath = fmt.Sprintf("%s/%s", downloadPath, path)
		} else {
			fullPath = downloadPath
		}
		if base.Exists(fullPath) && t.FileComplete(path) {
			fp := strings.Split(path, "/")
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fp[len(fp)-1]))
			c.File(fullPath)
		}
	} else {
		c.JSON(http.StatusNotFound, DataJson(false, "NotFound"))
	}
}

func torrentDownload(c *gin.Context) {
	hash := c.Param("hash")
	mi, err := db.SelectMetaInfo(hash)
	if err != nil {
		c.JSON(http.StatusOK, DataJson(false, err.Error()))
		return
	}
	if b, err := bencode.Marshal(mi); err == nil {
		c.Data(http.StatusOK, "application/x-bittorrent", b)
	} else {
		c.JSON(http.StatusOK, DataJson(false, err.Error()))
	}
}

func InitTorrentRouter(groupRouter *gin.RouterGroup) {
	groupRouter.GET("/info", torrentInfo)
	groupRouter.POST("/upload", torrentUpload)
	groupRouter.GET("/file/download", torrentFileDownload)
	groupRouter.GET("/download/:hash", torrentDownload)
}
