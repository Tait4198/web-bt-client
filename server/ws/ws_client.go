package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait = 30 * time.Second

	pongWait = 30 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 1024
)

type WebSocketClient struct {
	id           string
	conn         *websocket.Conn
	wsm          *WebSocketManager
	readMessage  chan []byte
	writeMessage chan interface{}
}

func (wsc *WebSocketClient) Send(msg interface{}) {
	wsc.writeMessage <- msg
}

func (wsc *WebSocketClient) readPump() {
	defer func() {
		wsc.wsm.clientClose <- wsc.id
		if err := wsc.conn.Close(); err != nil {
			log.Printf("Client %s Close Error %v\n", wsc.id, err)
		}
	}()
	wsc.conn.SetReadLimit(maxMessageSize)
	if err := wsc.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("Client %s SetReadDeadline Error %v\n", wsc.id, err)
	}
	wsc.conn.SetPongHandler(func(string) error {
		if err := wsc.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Printf("Client %s SetReadDeadline Error %v\n", wsc.id, err)
		}
		return nil
	})
	for {
		_, message, err := wsc.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Client %s readPump error: %v", wsc.id, err)
			}
			break
		}
		wsc.readMessage <- message
	}
}

func (wsc *WebSocketClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := wsc.conn.Close(); err != nil {
			log.Printf("Client %s Close Error %v\n", wsc.id, err)
		}
	}()
	for {
		select {
		case message, ok := <-wsc.writeMessage:
			if !ok {
				if err := wsc.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("Client %s Send Close Error %v\n", wsc.id, err)
				}
				return
			}
			if err := wsc.conn.WriteJSON(message); err != nil {
				log.Printf("Client %s Send Message Error %v\n", wsc.id, err)
			}
		case <-ticker.C:
			if err := wsc.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("Client %s WriteDeadline Error %v\n", wsc.id, err)
				return
			}
			if err := wsc.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Client %s Send Ping Error %v\n", wsc.id, err)
				return
			}
		}
	}
}

func newWebSocketClient(id string, conn *websocket.Conn, wsm *WebSocketManager) *WebSocketClient {
	wsc := &WebSocketClient{
		id:           id,
		conn:         conn,
		wsm:          wsm,
		readMessage:  make(chan []byte, 16),
		writeMessage: make(chan interface{}, 16),
	}
	go wsc.readPump()
	go wsc.writePump()
	return wsc
}
