package bencode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	val := "abc"
	buf := make([]byte, 10, 10)
	wLen := EncodeString(buf, val)
	str, _ := DecodeString(buf[:wLen+1])
	assert.Equal(t, val, str)

	val = ""
	for i := 0; i < 20; i++ {
		val += string(byte('a' + i))
	}
	buf = make([]byte, 100, 100)
	wLen = EncodeString(buf, val)
	str, _ = DecodeString(buf[:wLen+1])
	assert.Equal(t, val, str)
}

func TestInt(t *testing.T) {
	val := 999
	buf := make([]byte, 10, 10)
	wLen := EncodeInt(buf, val)
	iv, _ := DecodeInt(buf[:wLen+1])
	assert.Equal(t, val, iv)

	val = -99
	buf = make([]byte, 10, 10)
	wLen = EncodeInt(buf, val)
	iv, _ = DecodeInt(buf[:wLen+1])
	assert.Equal(t, val, iv)
}
