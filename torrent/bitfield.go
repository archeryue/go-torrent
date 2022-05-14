package torrent

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
