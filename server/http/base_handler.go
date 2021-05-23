package http

import "github.com/gin-gonic/gin"

func getPath(context *gin.Context) {
	// todo 路径信息
}

func InitBaseRouter(groupRouter *gin.RouterGroup) {
	groupRouter.POST("/path", getPath)
}
