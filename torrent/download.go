package torrent

type TorrentTask struct {
	PeerId   [20]byte
	PeerList []PeerInfo
	InfoSHA  [SHALEN]byte
	FileName string
	FileLen  int
	PieceLen int
	PieceSHA [][SHALEN]byte
}

func Download(task *TorrentTask) error {
	//TODO: check local tmp file
	//TODO: download piceces and check
	//TODO: write picece bytes into local tmp file
	//TODO: change tmp file name
	return nil
}
