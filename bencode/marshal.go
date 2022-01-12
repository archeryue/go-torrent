package bencode

import (
	"errors"
	"io"
	"reflect"
	"strings"
)

func Unmarshal(r io.Reader, s interface{}) error {
	o, err := Parse(r)
	if err != nil {
		return err
	}
	dict, err := o.Dict()
	if err != nil {
		return errors.New("src must be dict")
	}
	p := reflect.ValueOf(s)
	if p.Kind() != reflect.Ptr {
		return errors.New("dest must be pointer")
	}
	p = p.Elem()
	return unmarshal(p, dict)
}

func unmarshal(p reflect.Value, dict map[string]*BObject) error {
	for i, n := 0, p.NumField(); i < n; i++ {
		fv := p.Field(i)
		if !fv.CanSet() {
			continue
		}
		ft := p.Type().Field(i)
		key := ft.Tag.Get("bencode")
		if len(key) == 0 {
			key = strings.ToLower(ft.Name)
		}
		fo := dict[key]
		if fo == nil {
			continue
		}
		switch fo.type_ {
		case BSTR:
			if ft.Type.Kind() != reflect.String {
				break
			}
			val, _ := fo.Str()
			fv.SetString(val)
		case BINT:
			if ft.Type.Kind() != reflect.Int {
				break
			}
			val, _ := fo.Int()
			fv.SetInt(int64(val))
		case BLIST:
			if ft.Type.Kind() != reflect.Ptr {
				break
			}
			//TODO: handle list
		case BDICT:
			//TODO: handle dict
		}
	}
	return nil
}

func Marshal(w io.Writer, s interface{}) int {
	//TODO: Marshal to Writer
	return 0
}
