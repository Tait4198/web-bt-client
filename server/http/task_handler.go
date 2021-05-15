package http

import (
	"github.com/gin-gonic/gin"
	"github.com/web-bt-client/task"
	"net/http"
)

type NewTaskParam struct {
	Uri string `json:"uri"`
}

func newUriTask(c *gin.Context) {
	newTaskParam := NewTaskParam{}
	err := c.BindJSON(&newTaskParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageJson(false, "无效参数"))
		return
	}
	tm := task.GetTaskManager()
	hash, err := tm.AddUriTask(newTaskParam.Uri)
	if err == nil {
		c.JSON(http.StatusOK, DataJson(true, hash))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func newFileTask(c *gin.Context) {

}

func pauseTask(c *gin.Context) {
	hash := c.DefaultQuery("hash", "")
	tm := task.GetTaskManager()
	err := tm.Stop(hash)
	if err == nil {
		c.JSON(http.StatusOK, DataJson(true, hash))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func resumeTask(c *gin.Context) {
	hash := c.DefaultQuery("hash", "")
	tm := task.GetTaskManager()
	err := tm.Start(hash)
	if err == nil {
		c.JSON(http.StatusOK, DataJson(true, hash))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func InitTaskRouter(groupRouter *gin.RouterGroup) {
	groupRouter.POST("/new/uri", newUriTask)
	groupRouter.POST("/new/file", newFileTask)
	groupRouter.GET("/pause", pauseTask)
	groupRouter.GET("/resume", resumeTask)
}
