package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/archeryue/go-torrent/bencode"
)

type rawInfo struct {
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

type rawFile struct {
	Announce string  `bencode:"announce"`
	Info     rawInfo `bencode:"info"`
}

const SHALEN int = 20

type TorrentFile struct {
	Announce string
	InfoSHA  [SHALEN]byte
	FileName string
	FileLen  int
	PieceLen int
	PieceSHA [][SHALEN]byte
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
	ret.FileName = raw.Info.Name
	ret.FileLen = raw.Info.Length
	ret.PieceLen = raw.Info.PieceLength

	// calculate info SHA
	buf := new(bytes.Buffer)
	wlen := bencode.Marshal(buf, raw.Info)
	if wlen == 0 {
		fmt.Println("raw file info error")
	}
	ret.InfoSHA = sha1.Sum(buf.Bytes())

	// calculate pieces SHA
	bys := []byte(raw.Info.Pieces)
	cnt := len(bys) / SHALEN
	hashes := make([][SHALEN]byte, cnt)
	for i := 0; i < cnt; i++ {
		copy(hashes[i][:], bys[i*SHALEN:(i+1)*SHALEN])
	}
	ret.PieceSHA = hashes
	return ret, nil
}
