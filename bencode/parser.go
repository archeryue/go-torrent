package bencode

import (
	"bufio"
	"io"
)

func Parse(r io.Reader) (*BObject, error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	//recursive descent parsing
	b, err := br.ReadByte()
	if err != nil {
		return nil, err
	}
	br.UnreadByte()
	var ret BObject
	switch {
	case b >= '0' && b <= '9':
		val, err := DecodeString(br)
		if err != nil {
			return nil, err
		}
		ret.type_ = BSTR
		ret.val_ = val
	case b == 'i':
		val, err := DecodeInt(br)
		if err != nil {
			return nil, err
		}
		ret.type_ = BINT
		ret.val_ = val
	case b == 'l':
		br.ReadByte()
		var list []*BObject
		for {
			first, err := br.ReadByte()
			if err != nil {
				return nil, err
			}
			if first == 'e' {
				break
			}
			br.UnreadByte()
			elem, err := Parse(br)
			if err != nil {
				return nil, err
			}
			list = append(list, elem)
		}
		ret.type_ = BLIST
		ret.val_ = list
	case b == 'd':
		br.ReadByte()
		dict := make(map[string]*BObject)
		for {
			first, err := br.ReadByte()
			if err != nil {
				return nil, err
			}
			if first == 'e' {
				break
			}
			br.UnreadByte()
			key, err := DecodeString(br)
			if err != nil {
				return nil, err
			}
			val, err := Parse(br)
			if err != nil {
				return nil, err
			}
			dict[key] = val
		}
		ret.type_ = BDICT
		ret.val_ = dict
	}
	return &ret, nil
}
