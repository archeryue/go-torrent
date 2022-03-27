package torrent

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFile(t *testing.T) {
	file, err := os.Open("../testfile/debian-iso.torrent")
	assert.Equal(t, nil, err)
	tf, err := ParseFile(bufio.NewReader(file))
	assert.Equal(t, nil, err)
	assert.Equal(t, "http://bttracker.debian.org:6969/announce", tf.Announce)
	assert.Equal(t, "debian-11.2.0-amd64-netinst.iso", tf.FileName)
	assert.Equal(t, 396361728, tf.FileLen)
	assert.Equal(t, 262144, tf.PieceLen)
	assert.Equal(t, 1512, len(tf.PieceSHA))
	var expectHASH = [20]byte{0x28, 0xc5, 0x51, 0x96, 0xf5, 0x77, 0x53, 0xc4, 0xa,
		0xce, 0xb6, 0xfb, 0x58, 0x61, 0x7e, 0x69, 0x95, 0xa7, 0xed, 0xdb}
	assert.Equal(t, expectHASH, tf.InfoSHA)
}
