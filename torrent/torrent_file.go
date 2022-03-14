package torrent

import (
	"io"
)

type TorrentFile struct {
	Announce string
	InfoSHA  [20]byte
	FileName string
	FileLen  int
	PieceLen int
	PieceSHA [][20]byte
}

func ParseFile(r io.Reader) (*TorrentFile, error) {
	return nil, nil
	//TODO
}
