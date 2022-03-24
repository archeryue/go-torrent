package torrent

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFile(t *testing.T) {
	file, err := os.Open("file/debian-iso.torrent")
	assert.Equal(t, nil, err)
	tf, err := ParseFile(bufio.NewReader(file))
	assert.Equal(t, nil, err)
	assert.Equal(t, "http://bttracker.debian.org:6969/announce", tf.Announce)
	assert.Equal(t, "debian-11.2.0-amd64-netinst.iso", tf.FileName)
	assert.Equal(t, 396361728, tf.FileLen)
	assert.Equal(t, 262144, tf.PieceLen)
	assert.Equal(t, 1512, len(tf.PieceSHA))
	var expectHASH = [20]byte{113, 156, 7, 79, 121, 140, 87, 203, 206, 138, 146, 212, 155, 195, 177, 215, 88, 122, 189, 170}
	assert.Equal(t, expectHASH, tf.InfoSHA)
}
