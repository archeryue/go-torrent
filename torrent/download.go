package torrent

func Download(tf *TorrentFile, peerId [20]byte, peers []PeerInfo) error {
	//TODO: check local tmp file
	//TODO: download piceces and check
	//TODO: write picece bytes into local tmp file
	return nil
}

func MakeFile(tf *TorrentFile) {
	//TODO: assemble tmp to file
}
