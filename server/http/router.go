package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/web-bt-client/base"
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
	if base.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	if !base.Release {
		router.Use(cors.Default())
	}

	router.MaxMultipartMemory = 8 << 20

	taskRouter := router.Group("task")
	InitTaskRouter(taskRouter)

	torrentRouter := router.Group("torrent")
	InitTorrentRouter(torrentRouter)

	wsRouter := router.Group("ws")
	InitWsRouter(wsRouter)

	baseRouter := router.Group("base")
	InitBaseRouter(baseRouter)

	router.Static("/web", "./html")
	return router
}
