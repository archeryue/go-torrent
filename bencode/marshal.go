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
	p := reflect.ValueOf(s)
	switch o.type_ {
	case BLIST:
		list, _ := o.List()
		unmarshalList(p, list)
	case BDICT:
		dict, _ := o.Dict()
		unmarshalDict(p, dict)
	default:
		return errors.New("src code must be struct or slice")
	}
	return nil
}

// p.Kind must be Ptr && p.Elem().Type().Kind() must be Slice
func unmarshalList(p reflect.Value, list []*BObject) error {
	if p.Kind() != reflect.Ptr || p.Elem().Type().Kind() != reflect.Slice {
		return errors.New("dest must be pointer to slice")
	}
	p = p.Elem()
	if len(list) == 0 {
		return nil
	}
	switch list[0].type_ {
	case BSTR:
		for _, o := range list {
			val, err := o.Str()
			if err != nil {
				return err
			}
			p = reflect.Append(p, reflect.ValueOf(val))
		}
	case BINT:
		for _, o := range list {
			val, err := o.Int()
			if err != nil {
				return err
			}
			p = reflect.Append(p, reflect.ValueOf(val))
		}
	case BLIST:
		for _, o := range list {
			val, err := o.List()
			if err != nil {
				return err
			}
			if p.Type().Elem().Kind() != reflect.Slice {
				return ErrTyp
			}
			lp := reflect.New(p.Type().Elem())
			err = unmarshalList(lp, val)
			if err != nil {
				return err
			}
			p = reflect.Append(p, lp.Elem())
		}
	case BDICT:
		for _, o := range list {
			val, err := o.Dict()
			if err != nil {
				return err
			}
			if p.Type().Elem().Kind() != reflect.Struct {
				return ErrTyp
			}
			dp := reflect.New(p.Type().Elem())
			err = unmarshalDict(dp, val)
			if err != nil {
				return err
			}
			p = reflect.Append(p, dp.Elem())
		}
	}
	return nil
}

// p.Kind() must be Ptr && p.Elem().Type().Kind() must be Struct
func unmarshalDict(p reflect.Value, dict map[string]*BObject) error {
	if p.Kind() != reflect.Ptr || p.Elem().Type().Kind() != reflect.Struct {
		return errors.New("dest must be pointer")
	}
	p = p.Elem()
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
			if ft.Type.Kind() != reflect.Slice {
				break
			}
			valp := reflect.New(ft.Type)
			list, _ := fo.List()
			unmarshalList(valp, list)
			fv.Set(valp.Elem())
		case BDICT:
			if ft.Type.Kind() != reflect.Struct {
				break
			}
			valp := reflect.New(ft.Type)
			dict, _ := fo.Dict()
			unmarshalDict(valp, dict)
			fv.Set(valp.Elem())
		}
	}
	return nil
}

func Marshal(w io.Writer, s interface{}) int {
	//TODO: Marshal to Writer
	return 0
}
