package task

import "github.com/web-bt-client/ws"

func BroadcastTaskStatus(task *TorrentTask, messageType MessageType, status bool) {
	ws.GetWebSocketManager().Broadcast(TorrentTaskStatus{
		TorrentBase: TorrentBase{
			InfoHash: task.param.InfoHash,
			Type:     messageType,
		},
		Status: status,
	})
}
