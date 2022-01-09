package bencode

import (
	"bufio"
	"errors"
	"io"
)

var (
	ErrNum = errors.New("expect num")
	ErrCol = errors.New("expect colon")
	ErrEpI = errors.New("expect i")
	ErrEpE = errors.New("expect e")
)

func checkNum(data byte) bool {
	return data >= '0' && data <= '9'
}

func readDecimal(r *bufio.Reader) (val int, len int) {
	sign := 1
	b, _ := r.ReadByte()
	len++
	if b == '-' {
		sign = -1
		b, _ = r.ReadByte()
		len++
	}
	for {
		if !checkNum(b) {
			r.UnreadByte()
			len--
			return sign * val, len
		}
		val = val*10 + int(b-'0')
		b, _ = r.ReadByte()
		len++
	}
}

func writeDecimal(w io.Writer, val int) (len int) {
	if val == 0 {
		w.Write([]byte{'0'})
		len++
		return
	}
	if val < 0 {
		w.Write([]byte{'-'})
		len++
		val *= -1
	}

	dividend := 1
	for {
		if dividend > val {
			dividend /= 10
			break
		}
		dividend *= 10
	}
	for {
		num := byte(val / dividend)
		w.Write([]byte{'0' + num})
		len++
		if dividend == 1 {
			return
		}
		val %= dividend
		dividend /= 10
	}
}

func EncodeString(w io.Writer, val string) int {
	strLen := len(val)
	wLen := writeDecimal(w, strLen)
	w.Write([]byte{':'})
	wLen++
	w.Write([]byte(val))
	wLen += strLen
	return wLen
}

func DecodeString(r io.Reader) (val string, err error) {
	br := bufio.NewReader(r)
	num, len := readDecimal(br)
	if len == 0 {
		return val, ErrNum
	}
	b, err := br.ReadByte()
	if b != ':' {
		return val, ErrCol
	}
	buf, err := br.Peek(num)
	val = string(buf)
	return
}

func EncodeInt(w io.Writer, val int) int {
	wLen := 0
	w.Write([]byte{'i'})
	wLen++
	nLen := writeDecimal(w, val)
	wLen += nLen
	w.Write([]byte{'e'})
	wLen++
	return wLen
}

func DecodeInt(r io.Reader) (val int, err error) {
	br := bufio.NewReader(r)
	b, err := br.ReadByte()
	if b != 'i' {
		return val, ErrEpI
	}
	val, _ = readDecimal(br)
	b, err = br.ReadByte()
	if b != 'e' {
		return val, ErrEpE
	}
	return
}
