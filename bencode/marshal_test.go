package bencode

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string `bencode:"name"`
	Age  int    `bencode:"age"`
}

type Role struct {
	Id   int
	User `bencode:"user"`
}

type Score struct {
	User  `bencode:"user"`
	Value []int `bencode:"value"`
}

type Team struct {
	Name   string `bencode:"name"`
	Size   int    `bencode:"size"`
	Member []User `bencode:"member"`
}

func TestMarshalBasic(t *testing.T) {
	buf := new(bytes.Buffer)
	str := "abc"
	len := Marshal(buf, str)
	assert.Equal(t, 5, len)
	assert.Equal(t, "3:abc", buf.String())

	buf.Reset()
	val := 199
	len = Marshal(buf, val)
	assert.Equal(t, 5, len)
	assert.Equal(t, "i199e", buf.String())
}

func TestUnmarshalList(t *testing.T) {
	str := "li85ei90ei95ee"
	l := &[]int{}
	Unmarshal(bytes.NewBufferString(str), l)
	assert.Equal(t, []int{85, 90, 95}, *l)

	buf := new(bytes.Buffer)
	length := Marshal(buf, l)
	assert.Equal(t, len(str), length)
	assert.Equal(t, str, buf.String())
}

func TestUnmarshalUser(t *testing.T) {
	str := "d4:name6:archer3:agei29ee"
	u := &User{}
	Unmarshal(bytes.NewBufferString(str), u)
	assert.Equal(t, "archer", u.Name)
	assert.Equal(t, 29, u.Age)

	buf := new(bytes.Buffer)
	length := Marshal(buf, u)
	assert.Equal(t, len(str), length)
	assert.Equal(t, str, buf.String())
}

func TestUnmarshalRole(t *testing.T) {
	str := "d2:idi1e4:userd4:name6:archer3:agei29eee"
	r := &Role{}
	Unmarshal(bytes.NewBufferString(str), r)
	assert.Equal(t, 1, r.Id)
	assert.Equal(t, "archer", r.Name)
	assert.Equal(t, 29, r.Age)

	buf := new(bytes.Buffer)
	length := Marshal(buf, r)
	assert.Equal(t, len(str), length)
	assert.Equal(t, str, buf.String())
}

func TestUnmarshalScore(t *testing.T) {
	str := "d4:userd4:name6:archer3:agei29ee5:valueli80ei85ei90eee"
	s := &Score{}
	Unmarshal(bytes.NewBufferString(str), s)
	assert.Equal(t, "archer", s.Name)
	assert.Equal(t, 29, s.Age)
	assert.Equal(t, []int{80, 85, 90}, s.Value)

	buf := new(bytes.Buffer)
	length := Marshal(buf, s)
	assert.Equal(t, len(str), length)
	assert.Equal(t, str, buf.String())
}

func TestUnmarshalTeam(t *testing.T) {
	str := "d4:name3:ace4:sizei2e6:memberld4:name6:archer3:agei29eed4:name5:nancy3:agei31eeee"
	team := &Team{}
	Unmarshal(bytes.NewBufferString(str), team)
	assert.Equal(t, "ace", team.Name)
	assert.Equal(t, 2, team.Size)

	buf := new(bytes.Buffer)
	length := Marshal(buf, team)
	assert.Equal(t, len(str), length)
	assert.Equal(t, str, buf.String())
}
