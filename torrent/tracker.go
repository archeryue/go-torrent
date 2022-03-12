package torrent

import (
	"net"
)

type PeerInfo struct {
	Ip net.IP
	Port uint16
}

func FindPeers(tfile *TorrentFile) []PeerInfo {
	return nil
}
