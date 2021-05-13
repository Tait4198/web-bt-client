package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type WebSocketManager struct {
	clientClose chan string
	clientMap   *ClientMap
}

func (wsm *WebSocketManager) NewWebSocketConn(conn *websocket.Conn) {
	clientId := uuid.New().String()
	client := newWebSocketClient(clientId, conn, wsm)
	wsm.clientMap.Store(clientId, client)
	log.Printf("Client Size %d ->  %s Connection \n", wsm.clientMap.Count(), clientId)
}

func (wsm *WebSocketManager) Broadcast(msg interface{}) {
	wsm.clientMap.Range(func(id string, wsc *WebSocketClient) {
		wsc.Send(msg)
	})
}

func (wsm *WebSocketManager) Run() {
	for {
		select {
		case clientId := <-wsm.clientClose:
			if wsm.clientMap.HasMember(clientId) {
				wsm.clientMap.Delete(clientId)
				log.Printf("Client Size %d ->  %s Remove \n", wsm.clientMap.Count(), clientId)
			}
		}
	}
}

var taskManager *WebSocketManager
var tmOnce sync.Once

func GetWebSocketManager() *WebSocketManager {
	tmOnce.Do(func() {
		taskManager = &WebSocketManager{
			clientClose: make(chan string),
			clientMap:   NewClientMap(),
		}
		go taskManager.Run()
	})
	return taskManager
}
