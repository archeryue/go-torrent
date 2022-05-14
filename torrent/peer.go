package torrent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

type MsgId uint8

const (
	MsgChoke       MsgId = 0
	MsgUnchoke     MsgId = 1
	MsgInterested  MsgId = 2
	MsgNotInterest MsgId = 3
	MsgHave        MsgId = 4
	MsgBitfield    MsgId = 5
	MsgRequest     MsgId = 6
	MsgPiece       MsgId = 7
	MsgCancel      MsgId = 8
)

type PeerMsg struct {
	Id      MsgId
	Payload []byte
}

type PeerConn struct {
	net.Conn
	Choked  bool
	Field   Bitfield
	peer    PeerInfo
	peerId  [IDLEN]byte
	infoSHA [SHALEN]byte
}

func handshake(conn net.Conn, infoSHA [SHALEN]byte, peerId [IDLEN]byte) error {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})
	// send HandshakeMsg
	req := NewHandShakeMsg(infoSHA, peerId)
	_, err := WriteHandShake(conn, req)
	if err != nil {
		fmt.Println("send handshake failed")
		return err
	}
	// read HandshakeMsg
	res, err := ReadHandshake(conn)
	if err != nil {
		fmt.Println("read handshake failed")
		return err
	}
	// check HandshakeMsg
	if !bytes.Equal(res.InfoSHA[:], infoSHA[:]) {
		fmt.Println("check handshake failed")
		return fmt.Errorf("handshake msg error: " + string(res.InfoSHA[:]))
	}
	return nil
}

func fillBitfield(c *PeerConn) error {
	c.SetDeadline(time.Now().Add(3 * time.Second))
	defer c.SetDeadline(time.Time{})
	msg, err := c.ReadMsg()
	if err != nil {
		fmt.Println("read peer msg failed")
		return err
	}
	if msg == nil || msg.Id != MsgBitfield {
		fmt.Println("peer msg type error")
		return fmt.Errorf("expect bitfield")
	}
	c.Field = msg.Payload
	return nil
}

func (c *PeerConn) ReadMsg() (*PeerMsg, error) {
	// read msg length
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(c, lenBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenBuf)
	// keep alive msg
	if length == 0 {
		return nil, nil
	}
	// read msg body
	msgBuf := make([]byte, length)
	_, err = io.ReadFull(c, lenBuf)
	if err != nil {
		return nil, err
	}
	return &PeerMsg{
		Id:      MsgId(msgBuf[0]),
		Payload: msgBuf[1:],
	}, nil
}

const LenBytes uint32 = 4

func (c *PeerConn) WriteMsg(m *PeerMsg) (int, error) {
	var buf []byte
	if m == nil {
		buf = make([]byte, LenBytes)
	}
	length := uint32(len(m.Payload) + 1) // +1 for id
	buf = make([]byte, LenBytes+length)
	binary.BigEndian.PutUint32(buf[0:LenBytes], length)
	buf[LenBytes] = byte(m.Id)
	copy(buf[LenBytes+1:], m.Payload)
	return c.Write(buf)
}

func NewConn(peer PeerInfo, infoSHA [SHALEN]byte, peerId [IDLEN]byte) (*PeerConn, error) {
	// setup tcp conn
	addr := net.JoinHostPort(peer.Ip.String(), strconv.Itoa(int(peer.Port)))
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		fmt.Println("set tcp conn failed: " + addr)
		return nil, err
	}
	// torrent p2p handshake
	err = handshake(conn, infoSHA, peerId)
	if err != nil {
		fmt.Println("handshake failed")
		return nil, err
	}
	c := &PeerConn{
		Conn:    conn,
		Choked:  true,
		peer:    peer,
		peerId:  peerId,
		infoSHA: infoSHA,
	}
	// fill bitfield
	err = fillBitfield(c)
	if err != nil {
		fmt.Println("fill bitfield failed")
		return nil, err
	}
	return c, nil
}
