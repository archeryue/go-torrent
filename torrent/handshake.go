package torrent

import (
	"fmt"
	"io"
)

type HandshakeMsg struct {
	PreStr  string
	InfoSHA [20]byte
	PeerId  [20]byte
}

func NewHandShakeMsg(infoSHA, peerId [20]byte) *HandshakeMsg {
	return &HandshakeMsg{
		PreStr:  "BitTorrent protocol",
		InfoSHA: infoSHA,
		PeerId:  peerId,
	}
}

func WriteHandShake(w io.Writer, msg *HandshakeMsg) (int, error) {
	buf := make([]byte, len(msg.PreStr)+49)
	buf[0] = byte(len(msg.PreStr))
	curr := 1
	curr += copy(buf[curr:], []byte(msg.PreStr))
	curr += copy(buf[curr:], make([]byte, 8)) // 8 reserved bytes
	curr += copy(buf[curr:], msg.InfoSHA[:])
	curr += copy(buf[curr:], msg.PeerId[:])
	return w.Write(buf)
}

func ReadHandshake(r io.Reader) (*HandshakeMsg, error) {
	lenBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lenBuf)
	if err != nil {
		return nil, err
	}
	prelen := int(lenBuf[0])

	if prelen == 0 {
		err := fmt.Errorf("prelen cannot be 0")
		return nil, err
	}

	msgBuf := make([]byte, 48+prelen)
	_, err = io.ReadFull(r, msgBuf)
	if err != nil {
		return nil, err
	}

	var peerId [IDLEN]byte
	var infoSHA [SHALEN]byte

	copy(infoSHA[:], msgBuf[prelen+8:prelen+8+20])
	copy(peerId[:], msgBuf[prelen+8+20:])

	return &HandshakeMsg{
		PreStr:  string(msgBuf[0:prelen]),
		InfoSHA: infoSHA,
		PeerId:  peerId,
	}, nil
}
