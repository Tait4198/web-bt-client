package task

import "fmt"

func (dt *TorrentTask) GetTorrentStats(includeChunks bool, includePeers bool) TorrentStatsWrapper {
	torrentStats := dt.torrent.Stats()

	torrentStatsWrapper := TorrentStatsWrapper{
		TorrentBase: TorrentBase{
			InfoHash: dt.torrent.InfoHash().String(),
			Type:     Stats,
		},
		TorrentDownload: TorrentDownload{
			Length:         dt.torrent.Length(),
			DownloadLength: dt.download.downloadLength,
			BytesCompleted: dt.torrent.BytesCompleted(),
		},
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

	if includeChunks {
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
	}

	if includePeers {
		torrentStatsWrapper.TotalPeers = torrentStats.TotalPeers
		torrentStatsWrapper.ActivePeers = torrentStats.ActivePeers
		torrentStatsWrapper.HalfOpenPeers = torrentStats.HalfOpenPeers
		torrentStatsWrapper.PendingPeers = torrentStats.PendingPeers
	}

	return torrentStatsWrapper
}

func (dt *TorrentTask) GetTorrentInfo(includeFile bool) (TorrentInfoWrapper, error) {
	t := dt.torrent
	if t.Info() == nil {
		return TorrentInfoWrapper{}, fmt.Errorf(" %s 未获取 MateInfo", t.InfoHash().String())
	}
	torrentInfoWrapper := TorrentInfoWrapper{
		TorrentBase: TorrentBase{
			InfoHash: dt.torrent.InfoHash().String(),
			Type:     Info,
		},
		TorrentDownload: TorrentDownload{
			Length:         dt.torrent.Length(),
			DownloadLength: dt.download.downloadLength,
			BytesCompleted: dt.torrent.BytesCompleted(),
		},
		Name: t.Name(),
	}

	completedPieces := 0
	for _, psr := range t.PieceStateRuns() {
		if psr.Complete {
			completedPieces += psr.Length
		}
	}
	torrentInfoWrapper.CompletedPieces = completedPieces
	torrentInfoWrapper.Pieces = t.NumPieces()

	if includeFile {
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
	}

	return torrentInfoWrapper, nil
}
