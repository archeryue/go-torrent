package bencode

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseString(t *testing.T) {
	var o *BObject
	buf := bytes.NewBufferString("3:abc")
	o, _ = Parse(buf)
	assert.Equal(t, BSTR, o.type_)
	assert.Equal(t, "abc", o.val_.(string))
}

func TestParseInt(t *testing.T) {
	var o *BObject
	buf := bytes.NewBufferString("i123e")
	o, _ = Parse(buf)
	assert.Equal(t, BINT, o.type_)
	assert.Equal(t, 123, o.val_.(int))
}

func TestParseList(t *testing.T) {
	var o *BObject
	var list []*BObject
	buf := bytes.NewBufferString("li123e6:archeri789ee")
	o, _ = Parse(buf)
	assert.Equal(t, BLIST, o.type_)
	list = o.val_.([]*BObject) 
	assert.Equal(t, 3, len(list))
	assert.Equal(t, BINT, list[0].type_)
	assert.Equal(t, 123, list[0].val_.(int))
	assert.Equal(t, BSTR, list[1].type_)
	assert.Equal(t, "archer", list[1].val_.(string))
	assert.Equal(t, BINT, list[2].type_)
	assert.Equal(t, 789, list[2].val_.(int))
}

func TestParseMap(t *testing.T) {
	var o *BObject
	var dict map[string]*BObject
	buf := bytes.NewBufferString("d4:name6:archer3:agei29ee")
	o, _ = Parse(buf)
	assert.Equal(t, BDICT, o.type_)
	dict = o.val_.(map[string]*BObject)
	var val *BObject
	val = dict["name"]
	assert.Equal(t, BSTR, val.type_)
	assert.Equal(t, "archer", val.val_.(string))
	val = dict["age"]
	assert.Equal(t, BINT, val.type_)
	assert.Equal(t, 29, val.val_.(int))
}
