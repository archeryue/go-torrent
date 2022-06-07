package torrent

import (
	"fmt"
	"os"
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

func (task *TorrentTask) peerRountine(peer PeerInfo, taskQueue chan *pieceTask, resultQueue chan *pieceResult) {

}

func (t *TorrentTask) downloadPiece(conn *PeerConn, task *pieceTask) (*pieceResult, error) {
	return &pieceResult{}, nil
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
		taskQueue <- &pieceTask{index, sha, (begin-end)}
	}
	// init goroutines for each peer
	for _, peer := range task.PeerList {
		go task.peerRountine(peer, taskQueue, resultQueue)
	}
	// collect piece result
	buf := make([]byte, task.FileLen)
	count := 0
	for count < len(task.PieceSHA) {
		res := <- resultQueue
		begin, end := task.getPieceBounds(res.index)
		copy(buf[begin:end], res.data)
		count++
		// print progress
		percent := float64(count) / float64(len(task.PieceSHA)) * 100
		fmt.Printf("downloading, progress : (%0.2f%%)", percent)
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
