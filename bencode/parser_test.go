package bencode

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func objAssertStr(t *testing.T, expect string, o *BObject) {
	assert.Equal(t, BSTR, o.type_)
	str, err := o.Str()
	assert.Equal(t, nil, err)
	assert.Equal(t, expect, str)
}

func objAssertInt(t *testing.T, expect int, o *BObject) {
	assert.Equal(t, BINT, o.type_)
	val, err := o.Int()
	assert.Equal(t, nil, err)
	assert.Equal(t, expect, val)
}

func TestParseString(t *testing.T) {
	var o *BObject
	buf := bytes.NewBufferString("3:abc")
	o, _ = Parse(buf)
	objAssertStr(t, "abc", o)
}

func TestParseInt(t *testing.T) {
	var o *BObject
	buf := bytes.NewBufferString("i123e")
	o, _ = Parse(buf)
	objAssertInt(t, 123, o)
}

func TestParseList(t *testing.T) {
	var o *BObject
	var list []*BObject
	buf := bytes.NewBufferString("li123e6:archeri789ee")
	o, _ = Parse(buf)
	assert.Equal(t, BLIST, o.type_)
	list, err := o.List()
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(list))
	objAssertInt(t, 123, list[0])
	objAssertStr(t, "archer", list[1])
	objAssertInt(t, 789, list[2])
}

func TestParseMap(t *testing.T) {
	var o *BObject
	var dict map[string]*BObject
	buf := bytes.NewBufferString("d4:name6:archer3:agei29ee")
	o, _ = Parse(buf)
	assert.Equal(t, BDICT, o.type_)
	dict, err := o.Dict()
	assert.Equal(t, nil, err)
	objAssertStr(t, "archer", dict["name"])
	objAssertInt(t, 29, dict["age"])
}
