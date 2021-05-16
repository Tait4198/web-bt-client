package task

import "fmt"

type TorrentInfoWrapper struct {
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`

	Length         int64 `json:"length"`
	DownloadLength int64 `json:"download_length"`
	BytesCompleted int64 `json:"bytes_completed"`

	Pieces          int `json:"pieces"`
	CompletedPieces int `json:"completed_pieces"`

	Files []TorrentInfoFileWrapper `json:"files"`
}

type TorrentInfoFileWrapper struct {
	Path []string `json:"path"`

	Length         int64 `json:"length"`
	BytesCompleted int64 `json:"bytes_completed"`

	Pieces          int `json:"pieces"`
	CompletedPieces int `json:"completed_pieces"`
}

type TorrentStatsWrapper struct {
	InfoHash string `json:"info_hash"`

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

func (dt *TorrentTask) GetTorrentStats() TorrentStatsWrapper {
	torrentStats := dt.torrent.Stats()

	torrentStatsWrapper := TorrentStatsWrapper{
		InfoHash:      dt.torrent.InfoHash().String(),
		TotalPeers:    torrentStats.TotalPeers,
		ActivePeers:   torrentStats.ActivePeers,
		HalfOpenPeers: torrentStats.HalfOpenPeers,
		PendingPeers:  torrentStats.PendingPeers,
	}

	bytesRead := torrentStats.BytesRead
	bytesReadData := torrentStats.BytesReadData
	bytesReadUsefulData := torrentStats.BytesReadUsefulData
	torrentStatsWrapper.BytesRead = bytesRead.Int64()
	torrentStatsWrapper.BytesReadData = bytesReadData.Int64()
	torrentStatsWrapper.BytesReadUsefulData = bytesReadUsefulData.Int64()

	bytesWritten := torrentStats.BytesWritten
	bytesWrittenData := torrentStats.BytesWrittenData
	torrentStatsWrapper.BytesWritten = bytesWritten.Int64()
	torrentStatsWrapper.BytesWrittenData = bytesWrittenData.Int64()

	chunksRead := torrentStats.ChunksRead
	chunksReadUseful := torrentStats.ChunksReadUseful
	chunksReadWasted := torrentStats.ChunksReadWasted
	chunksWritten := torrentStats.ChunksWritten
	metadataChunksRead := torrentStats.MetadataChunksRead
	torrentStatsWrapper.ChunksRead = chunksRead.Int64()
	torrentStatsWrapper.ChunksReadUseful = chunksReadUseful.Int64()
	torrentStatsWrapper.ChunksReadWasted = chunksReadWasted.Int64()
	torrentStatsWrapper.ChunksWritten = chunksWritten.Int64()
	torrentStatsWrapper.MetadataChunksRead = metadataChunksRead.Int64()

	return torrentStatsWrapper
}

func (dt *TorrentTask) GetTorrentInfo() (TorrentInfoWrapper, error) {
	t := dt.torrent
	if t.Info() == nil {
		return TorrentInfoWrapper{}, fmt.Errorf(" %s 未获取 MateInfo", t.InfoHash().String())
	}
	torrentInfoWrapper := TorrentInfoWrapper{
		InfoHash:       t.InfoHash().String(),
		Name:           t.Name(),
		Length:         t.Length(),
		DownloadLength: dt.download.downloadLength,
		BytesCompleted: t.BytesCompleted(),
	}

	completedPieces := 0
	for _, psr := range t.PieceStateRuns() {
		if psr.Complete {
			completedPieces += psr.Length
		}
	}
	torrentInfoWrapper.CompletedPieces = completedPieces
	torrentInfoWrapper.Pieces = t.NumPieces()

	var files = make([]TorrentInfoFileWrapper, len(t.Files()))
	for i, f := range t.Files() {
		filePieces := 0
		fileCompletedPieces := 0
		for _, state := range f.State() {
			filePieces++
			if state.Complete {
				fileCompletedPieces++
			}
		}
		torrentFile := TorrentInfoFileWrapper{
			Path:            f.FileInfo().Path,
			Length:          f.Length(),
			BytesCompleted:  f.BytesCompleted(),
			Pieces:          filePieces,
			CompletedPieces: fileCompletedPieces,
		}
		files[i] = torrentFile
	}
	if len(files) == 1 && files[0].Path == nil {
		files[0].Path = []string{t.Name()}
	}

	torrentInfoWrapper.Files = files

	return torrentInfoWrapper, nil
}
