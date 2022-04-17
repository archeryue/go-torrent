package torrent

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestTracker(t *testing.T) {
	file, _ := os.Open("../testfile/debian-iso.torrent")
	tf, _ := ParseFile(bufio.NewReader(file))

	peers := FindPeers(tf)
	for i, p := range peers {
		fmt.Printf("Peer %d, Ip: %s, Port: %d\n", i, p.Ip, p.Port)
	}
}
