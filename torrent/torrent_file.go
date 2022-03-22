package torrent

import (
	"fmt"
	"io"

	"github.com/archeryue/go-torrent/bencode"
)

type rawInfo struct {
	Name string		`bencode:"name"`
	Length int		`bencode:"length"`
	pieceLength int	`bencode:"piece length"`
	pieces []byte	`bencode:"pieces"`
}

type rawFile struct {
	Announce string	`bencode:"announce"`
	info rawInfo	`bencode:"info"`
}

type TorrentFile struct {
	Announce string
	InfoSHA  [20]byte
	FileName string
	FileLen  int
	PieceLen int
	PieceSHA [][20]byte
}

func ParseFile(r io.Reader) (*TorrentFile, error) {
	raw := new(rawFile)
	err := bencode.Unmarshal(r, raw)
	if err != nil {
		fmt.Println("Fail to parse torrent file")
		return nil, err
	}
	ret := new(TorrentFile)
	ret.Announce = raw.Announce
	ret.FileName = raw.info.Name
	ret.FileLen = raw.info.Length
	ret.PieceLen = raw.info.pieceLength
	//TODO: InfoSHA & PieceSHA
	return ret, nil
}
