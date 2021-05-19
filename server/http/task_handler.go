package http

import (
	"github.com/gin-gonic/gin"
	"github.com/web-bt-client/task"
	"net/http"
)

type UriTaskParam struct {
	Uri       string     `json:"uri"`
	TaskParam task.Param `json:"task_param"`
}

func newUriTask(c *gin.Context) {
	uriTaskParam := UriTaskParam{}
	err := c.BindJSON(&uriTaskParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageJson(false, "无效参数"))
		return
	}
	tm := task.GetTaskManager()
	hash, err := tm.AddUriTask(uriTaskParam.Uri, uriTaskParam.TaskParam)
	if err == nil {
		c.JSON(http.StatusOK, DataJson(true, hash))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func newFileTask(c *gin.Context) {

}

func stopTask(c *gin.Context) {
	hash := c.DefaultQuery("hash", "")
	tm := task.GetTaskManager()
	err := tm.Stop(hash, true)
	if err == nil {
		c.JSON(http.StatusOK, DataJson(true, hash))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func startTask(c *gin.Context) {
	tp := task.Param{}
	err := c.BindJSON(&tp)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageJson(false, "无效参数"))
		return
	}
	tm := task.GetTaskManager()
	if err := tm.Start(tp, true); err == nil {
		c.JSON(http.StatusOK, DataJson(true, tp.InfoHash))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func restartTask(c *gin.Context) {
	tp := task.Param{}
	err := c.BindJSON(&tp)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageJson(false, "无效参数"))
		return
	}
	tm := task.GetTaskManager()
	if err := tm.Restart(tp); err == nil {
		c.JSON(http.StatusOK, DataJson(true, tp.InfoHash))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func taskList(c *gin.Context) {
	tasks, err := task.GetTaskManager().GetTasks()
	if err == nil {
		c.JSON(http.StatusOK, DataJson(true, tasks))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func InitTaskRouter(groupRouter *gin.RouterGroup) {
	groupRouter.POST("/new/uri", newUriTask)
	groupRouter.POST("/new/file", newFileTask)
	groupRouter.GET("/stop", stopTask)
	groupRouter.POST("/start", startTask)
	groupRouter.POST("/restart", restartTask)
	groupRouter.GET("/list", taskList)
}
