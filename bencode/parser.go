package bencode

import "errors"

var (
	ErrNum = errors.New("expect num")
	ErrCol = errors.New("expect colon")
	ErrEpI = errors.New("expect i")
	ErrEpE = errors.New("expect e")
)

func checkNum(data byte) bool {
	return data >= '0' && data <= '9'
}

func getDecimal(data []byte) (val int, len int) {
	sign := 1
	if data[0] == '-' {
		sign = -1
		len++
	}
	for {
		if !checkNum(data[len]) {
			return sign * val, len
		}
		val = val*10 + int(data[len]-'0')
		len++
	}
}

func setDecimal(buf []byte, val int) (len int) {
	if val == 0 {
		buf[len] = '0'; len++
		return
	}
	if val < 0 {
		buf[len] = '-'; len++
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
		buf[len] = '0' + num; len++
		if dividend == 1 {
			return
		}
		val %= dividend; dividend /= 10
	}
}

func WriteString(buf []byte, val string) int {
	strLen := len(val)
	wLen := setDecimal(buf, strLen)
	buf[wLen] = ':'
	wLen++
	copy(buf[wLen:], []byte(val))
	wLen += strLen
	return wLen
}

func ParseString(src []byte) (val string, err error) {
	num, len := getDecimal(src)
	if len == 0 {
		return val, ErrNum
	}
	if src[len] != ':' {
		return val, ErrCol
	}
	val = string(src[len+1 : len+1+num])
	return
}

func WriteInt(buf []byte, val int) int {
	wLen := 0
	buf[0] = 'i'
	wLen++
	nLen := setDecimal(buf[wLen:], val)
	wLen += nLen
	buf[wLen] = 'e'
	wLen++
	return wLen
}

func ParseInt(src []byte) (val int, err error) {
	if src[0] != 'i' {
		return val, ErrEpI
	}
	val, len := getDecimal(src[1:])
	if src[len+1] != 'e' {
		return val, ErrEpE
	}
	return
}
