package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/web-bt-client/ws"
	"net/http"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsConn(c *gin.Context) {
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, MessageJson(false, "连接失败"))
		return
	} else {
		ws.GetWebSocketManager().NewWebSocketConn(conn)
	}
}

func InitWsRouter(groupRouter *gin.RouterGroup) {
	groupRouter.GET("/conn", wsConn)
}
