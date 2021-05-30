package http

import (
	"github.com/gin-gonic/gin"
	"github.com/web-bt-client/task"
	"github.com/web-bt-client/ws"
	"net/http"
	"strconv"
)

func createTask(c *gin.Context) {
	taskParam := task.Param{}
	err := c.BindJSON(&taskParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageJson(false, "无效参数"))
		return
	}
	tm := task.GetTaskManager()
	hash, err := tm.CreateTask(taskParam)
	if err == nil {
		c.JSON(http.StatusOK, DataJson(true, hash))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func deleteTask(c *gin.Context) {
	hash := c.DefaultQuery("hash", "")
	tm := task.GetTaskManager()
	err := tm.Delete(hash)
	if err == nil {
		c.JSON(http.StatusOK, DataJson(true, hash))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
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

func sentStatus(c *gin.Context) {
	status := c.DefaultQuery("status", "1")
	typeStr := c.DefaultQuery("type", "")
	hash := c.DefaultQuery("hash", "")
	typeInt, _ := strconv.ParseInt(typeStr, 10, 32)
	statusBool := status == "1"
	ws.GetWebSocketManager().Broadcast(task.TorrentTaskStatus{
		TorrentBase: task.TorrentBase{
			InfoHash: hash,
			Type:     task.MessageType(typeInt),
		},
		Status: statusBool,
	})
}

func taskExists(c *gin.Context) {
	hash := c.DefaultQuery("hash", "")
	tm := task.GetTaskManager()
	c.JSON(http.StatusOK, DataJson(true, tm.TaskExists(hash)))
}

func InitTaskRouter(groupRouter *gin.RouterGroup) {
	groupRouter.POST("/create", createTask)
	groupRouter.GET("/delete", deleteTask)
	groupRouter.GET("/stop", stopTask)
	groupRouter.POST("/start", startTask)
	groupRouter.POST("/restart", restartTask)
	groupRouter.GET("/list", taskList)
	groupRouter.GET("/exists", taskExists)
	groupRouter.GET("/send", sentStatus)
}
