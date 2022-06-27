package torrent

import "strconv"

type Bitfield []byte

func (field Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	if byteIndex < 0 || byteIndex >= len(field) {
		return false
	}
	return field[byteIndex]>>uint(7-offset)&1 != 0
}

func (field Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8
	if byteIndex < 0 || byteIndex >= len(field) {
		return
	}
	field[byteIndex] |= 1 << uint(7-offset)
}

func (field Bitfield) String() string {
	str := "piece# "
	for i := 0; i < len(field)*8; i++ {
		if field.HasPiece(i) {
			str = str + strconv.Itoa(i) + " "
		}
	}
	return str
}
