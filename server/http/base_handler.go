package http

import (
	"github.com/gin-gonic/gin"
	"github.com/web-bt-client/base"
	"net/http"
)

func getPath(c *gin.Context) {
	parent := c.DefaultQuery("parent", "")
	if paths, err := base.GetPath(parent); err == nil {
		c.JSON(http.StatusOK, DataJson(true, paths))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func getSpace(c *gin.Context) {
	path := c.DefaultQuery("path", "")
	if paths, err := base.GetSpace(path); err == nil {
		c.JSON(http.StatusOK, DataJson(true, paths))
	} else {
		c.JSON(http.StatusOK, MessageJson(false, err.Error()))
	}
}

func InitBaseRouter(groupRouter *gin.RouterGroup) {
	groupRouter.GET("/path", getPath)
	groupRouter.GET("/space", getSpace)
}
