package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"
	"time"
)

type TorrentTask struct {
	PeerId   [20]byte
	PeerList []PeerInfo
	InfoSHA  [SHALEN]byte
	FileName string
	FileLen  int
	PieceLen int
	PieceSHA [][SHALEN]byte
}

type pieceTask struct {
	index  int
	sha    [SHALEN]byte
	length int
}

type taskState struct {
	index      int
	conn       *PeerConn
	requested  int
	downloaded int
	backlog    int
	data       []byte
}

type pieceResult struct {
	index int
	data  []byte
}

const BLOCKSIZE = 16384
const MAXBACKLOG = 5

func (state *taskState) handleMsg() error {
	msg, err := state.conn.ReadMsg()
	if err != nil {
		return err
	}
	// handle keep-alive
	if msg == nil {
		return nil
	}
	switch msg.Id {
	case MsgChoke:
		state.conn.Choked = true
	case MsgUnchoke:
		state.conn.Choked = false
	case MsgHave:
		index, err := GetHaveIndex(msg)
		if err != nil {
			return err
		}
		state.conn.Field.SetPiece(index)
	case MsgPiece:
		n, err := CopyPieceData(state.index, state.data, msg)
		if err != nil {
			return err
		}
		state.downloaded += n
		state.backlog--
	}
	return nil
}

func downloadPiece(conn *PeerConn, task *pieceTask) (*pieceResult, error) {
	state := &taskState{
		index: task.index,
		conn:  conn,
		data:  make([]byte, task.length),
	}
	conn.SetDeadline(time.Now().Add(15 * time.Second))
	defer conn.SetDeadline(time.Time{})

	for state.downloaded < task.length {
		if !conn.Choked {
			for state.backlog < MAXBACKLOG && state.requested < task.length {
				length := BLOCKSIZE
				if task.length-state.requested < length {
					length = task.length - state.requested
				}
				msg := NewRequestMsg(state.index, state.requested, length)
				_, err := state.conn.WriteMsg(msg)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += length
			}
		}
		err := state.handleMsg()
		if err != nil {
			return nil, err
		}
	}
	return &pieceResult{state.index, state.data}, nil
}

func checkPiece(task *pieceTask, res *pieceResult) bool {
	sha := sha1.Sum(res.data)
	if !bytes.Equal(task.sha[:], sha[:]) {
		fmt.Printf("check integrity failed, index :%v\n", res.index)
		return false
	}
	return true
}

func (t *TorrentTask) peerRoutine(peer PeerInfo, taskQueue chan *pieceTask, resultQueue chan *pieceResult) {
	// set up conn with peer
	conn, err := NewConn(peer, t.InfoSHA, t.PeerId)
	if err != nil {
		fmt.Println("fail to connect peer : " + peer.Ip.String())
		return
	}
	defer conn.Close()

	fmt.Println("complete handshake with peer : " + peer.Ip.String())
	conn.WriteMsg(&PeerMsg{MsgInterested, nil})
	// get piece task & download
	for task := range taskQueue {
		if !conn.Field.HasPiece(task.index) {
			taskQueue <- task
			continue
		}
		fmt.Printf("get task, index: %v, peer : %v\n", task.index, peer.Ip.String())
		res, err := downloadPiece(conn, task)
		if err != nil {
			taskQueue <- task
			fmt.Println("fail to download piece" + err.Error())
			return
		}
		if !checkPiece(task, res) {
			taskQueue <- task
			continue
		}
		resultQueue <- res
	}
}

func (t *TorrentTask) getPieceBounds(index int) (bengin, end int) {
	bengin = index * t.PieceLen
	end = bengin + t.PieceLen
	if end > t.FileLen {
		end = t.FileLen
	}
	return
}

func Download(task *TorrentTask) error {
	fmt.Println("start downloading " + task.FileName)
	// split pieceTasks and init task&result channel
	taskQueue := make(chan *pieceTask, len(task.PieceSHA))
	resultQueue := make(chan *pieceResult)
	for index, sha := range task.PieceSHA {
		begin, end := task.getPieceBounds(index)
		taskQueue <- &pieceTask{index, sha, (end - begin)}
	}
	// init goroutines for each peer
	for _, peer := range task.PeerList {
		go task.peerRoutine(peer, taskQueue, resultQueue)
	}
	// collect piece result
	buf := make([]byte, task.FileLen)
	count := 0
	for count < len(task.PieceSHA) {
		res := <-resultQueue
		begin, end := task.getPieceBounds(res.index)
		copy(buf[begin:end], res.data)
		count++
		// print progress
		percent := float64(count) / float64(len(task.PieceSHA)) * 100
		fmt.Printf("downloading, progress : (%0.2f%%)\n", percent)
	}
	close(taskQueue)
	close(resultQueue)
	// create file & copy data
	file, err := os.Create(task.FileName)
	if err != nil {
		fmt.Println("fail to create file: " + task.FileName)
		return err
	}
	_, err = file.Write(buf)
	if err != nil {
		fmt.Println("fail to write data")
		return err
	}
	return nil
}
