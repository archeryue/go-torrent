package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/archeryue/go-torrent/torrent"
)

func main() {
	//parse torrent file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("open file error")
		return
	}
	defer file.Close()
	tf, err := torrent.ParseFile(bufio.NewReader(file))
	if err != nil {
		fmt.Println("parse file error")
		return
	}
	// generate peerId
	var peerId [20]byte
	_, err = rand.Read(peerId[:])
	if err != nil {
		fmt.Println("generate peerId error")
		return
	}
	//connect tracker & find peers
	peers := torrent.FindPeers(tf, peerId)
	if len(peers) == 0 {
		fmt.Println("can not find peers")
	}
	//download from peers & make file
	torrent.Download(tf, peerId, peers)
	torrent.MakeFile(tf)
}
