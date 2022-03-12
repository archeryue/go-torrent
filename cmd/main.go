package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/archeryue/go-torrent/torrent"
)

func main() {
	//TODO: parse torrent file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("open file error")
	}
	tf, err := torrent.ParseFile(bufio.NewReader(file))
	if err != nil {
		fmt.Println("parse file error")
	}
	//TODO: connect tracker & find peers
	peers := torrent.FindPeers(tf)
	if len(peers) == 0 {
		fmt.Println("can not find peers")
	}
	//TODO: download from peers & make file
	torrent.Download()
	torrent.MakeFile()
}
