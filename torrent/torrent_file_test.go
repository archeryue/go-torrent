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
	var expectHASH = [20]byte{0x71, 0x2c, 0xea, 0x2f, 0x4b, 0xd8, 0x85, 0xa, 0xfa,
		0x19, 0xf9, 0x29, 0x45, 0xb0, 0xfa, 0xfe, 0x54, 0x97, 0xb9, 0x0e}
	assert.Equal(t, expectHASH, tf.InfoSHA)
}
