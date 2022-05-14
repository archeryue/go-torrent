package torrent

import "net"

type PeerConn struct {
	Conn   net.Conn
	Choked bool
	//bitfield
	peer    PeerInfo
	peerId  [IDLEN]byte
	infoSHA [SHALEN]byte
}

func New(peer PeerInfo, peerId [IDLEN]byte, infoSHA [SHALEN]byte) (conn *PeerConn, err error) {
	//TODO: set up conn
	return nil, nil
}

func (c *PeerConn) Handshake() error {
	return nil
}
