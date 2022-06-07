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

type pieceTask struct {
	index  int
	sha    [SHALEN]byte
	length int
}

type taskState struct {
	index      int
	conn       *PeerConn
	requested  int
	downloaded int
	backlog    int
	data       []byte
}

type pieceResult struct {
	index int
	data  []byte
}

const BLOCKSIZE = 16384
const MAXBACKLOG = 5

func downloadPiece(conn *PeerConn, task pieceTask) (pieceResult, error) {
	return pieceResult{}, nil
}

func Download(task *TorrentTask) error {
	// TODO: split pieceTasks and init task&result channel
	// TODO: init goroutines for each peer
	// TODO: check result channel
	// TODO: create file & copy data
	return nil
}
