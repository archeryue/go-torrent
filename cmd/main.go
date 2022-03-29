package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/archeryue/go-torrent/torrent"
)

func main() {
	//parse torrent file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("open file error")
	}
	defer file.Close()
	tf, err := torrent.ParseFile(bufio.NewReader(file))
	if err != nil {
		fmt.Println("parse file error")
	}
	//connect tracker & find peers
	peers := torrent.FindPeers(tf)
	if len(peers) == 0 {
		fmt.Println("can not find peers")
	}
	//download from peers & make file
	torrent.Download(tf, peers)
	torrent.MakeFile(tf)
}
