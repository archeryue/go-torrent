module archeryue/go-torrent/main

go 1.17

replace github.com/archeryue/go-torrent/torrent => ../torrent

replace github.com/archeryue/go-torrent/bencode => ../bencode

require github.com/archeryue/go-torrent/torrent v0.0.0-00010101000000-000000000000

require github.com/archeryue/go-torrent/bencode v0.0.0-20220320082858-9628ecdf6dfa // indirect
