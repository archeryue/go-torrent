package torrent

import (
	"net"
)

type PeerInfo struct {
	Ip   net.IP
	Port uint16
}

func FindPeers(tf *TorrentFile) []PeerInfo {
	return nil
}
