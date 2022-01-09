package bencode

import "io"

func Parse(r io.Reader) (*BObject, error) {
	//TODO: recursive descent parsing
	return &BObject{
		type_: BINT,
		val_: 0,
	}, nil
}
