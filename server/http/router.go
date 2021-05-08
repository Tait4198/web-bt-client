package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func MessageJson(status bool, message string) gin.H {
	return gin.H{
		"message": message,
		"status":  status,
	}
}

func DataJson(status bool, data interface{}) gin.H {
	return gin.H{
		"data":   data,
		"status": status,
	}
}

func Router() http.Handler {
	router := gin.New()
	taskRouter := router.Group("task")
	InitTaskRouter(taskRouter)
	return router
}
