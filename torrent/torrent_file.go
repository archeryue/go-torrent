package torrent

import (
	"io"
)

type TorrentFile struct {
    Announce    string
    InfoHash    [20]byte
    PieceHashes [][20]byte
    PieceLength int
    Length      int
    Name        string
	//TODO
}

func ParseFile(r io.Reader) (*TorrentFile, error) {
	return nil, nil
	//TODO
}
