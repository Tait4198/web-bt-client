package task

type MessageType int

const (
	TorrentStats MessageType = 1000
	TorrentInfo  MessageType = 1001
	TorrentWait  MessageType = 1002
)

type TorrentBase struct {
	InfoHash string      `json:"info_hash"`
	Type     MessageType `json:"type"`
}

type TorrentDownload struct {
	Length         int64 `json:"length"`
	DownloadLength int64 `json:"download_length"`
	BytesCompleted int64 `json:"bytes_completed"`
}

type TorrentInfoWrapper struct {
	TorrentBase
	TorrentDownload
	Name string `json:"name"`

	Pieces          int `json:"pieces"`
	CompletedPieces int `json:"completed_pieces"`

	Files []TorrentInfoFileWrapper `json:"files,omitempty"`
}

type TorrentInfoFileWrapper struct {
	Path []string `json:"path"`

	Length         int64 `json:"length"`
	BytesCompleted int64 `json:"bytes_completed"`

	Pieces          int `json:"pieces"`
	CompletedPieces int `json:"completed_pieces"`
}

type TorrentStatsWrapper struct {
	TorrentBase
	TorrentDownload

	BytesRead           int64 `json:"bytes_read"`
	BytesReadData       int64 `json:"bytes_read_data"`
	BytesReadUsefulData int64 `json:"bytes_read_useful_data"`

	BytesWritten     int64 `json:"bytes_written"`
	BytesWrittenData int64 `json:"bytes_written_data"`

	ChunksRead         int64 `json:"chunks_read"`
	ChunksReadUseful   int64 `json:"chunks_read_useful"`
	ChunksReadWasted   int64 `json:"chunks_read_wasted"`
	ChunksWritten      int64 `json:"chunks_written"`
	MetadataChunksRead int64 `json:"metadata_chunks_read"`

	TotalPeers    int `json:"total_peers"`
	ActivePeers   int `json:"active_peers"`
	HalfOpenPeers int `json:"half_open_peers"`
	PendingPeers  int `json:"pending_peers"`
}

type TorrentTaskWait struct {
	TorrentBase
	Status bool `json:"status"`
}
