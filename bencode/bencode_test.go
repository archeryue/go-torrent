package bencode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	val := "abc"
	buf := make([]byte, 10, 10)
	wLen := WriteString(buf, val)
	str, _ := ParseString(buf[:wLen+1])
	assert.Equal(t, val, str)
}

func TestInt(t *testing.T) {
	val := 999
	buf := make([]byte, 10, 10)
	wLen := WriteInt(buf, val)
	iv, _ := ParseInt(buf[:wLen+1])
	assert.Equal(t, val, iv)
}
