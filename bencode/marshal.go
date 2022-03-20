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
	if p.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer")
	}
	switch o.type_ {
	case BLIST:
		list, _ := o.List()
		l := reflect.MakeSlice(p.Elem().Type(), len(list), len(list))
		p.Elem().Set(l)
		err = unmarshalList(p, list)
		if err != nil {
			return err
		}
	case BDICT:
		dict, _ := o.Dict()
		err = unmarshalDict(p, dict)
		if err != nil {
			return err
		}
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
	vl := p.Elem()
	if len(list) == 0 {
		return nil
	}
	switch list[0].type_ {
	case BSTR:
		for i, o := range list {
			val, err := o.Str()
			if err != nil {
				return err
			}
			vl.Index(i).SetString(val)
		}
	case BINT:
		for i, o := range list {
			val, err := o.Int()
			if err != nil {
				return err
			}
			vl.Index(i).SetInt(int64(val))
		}
	case BLIST:
		for i, o := range list {
			val, err := o.List()
			if err != nil {
				return err
			}
			if vl.Type().Elem().Kind() != reflect.Slice {
				return ErrTyp
			}
			lp := reflect.New(vl.Type().Elem())
			ll := reflect.MakeSlice(vl.Type().Elem(), len(val), len(val))
			lp.Elem().Set(ll)
			err = unmarshalList(lp, val)
			if err != nil {
				return err
			}
			vl.Index(i).Set(lp.Elem())
		}
	case BDICT:
		for i, o := range list {
			val, err := o.Dict()
			if err != nil {
				return err
			}
			if vl.Type().Elem().Kind() != reflect.Struct {
				return ErrTyp
			}
			dp := reflect.New(vl.Type().Elem())
			err = unmarshalDict(dp, val)
			if err != nil {
				return err
			}
			vl.Index(i).Set(dp.Elem())
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
			list, _ := fo.List()
			valp := reflect.New(ft.Type)
			vall := reflect.MakeSlice(ft.Type, len(list), len(list))
			valp.Elem().Set(vall)
			err := unmarshalList(valp, list)
			if err != nil {
				break
			}
			fv.Set(valp.Elem())
		case BDICT:
			if ft.Type.Kind() != reflect.Struct {
				break
			}
			valp := reflect.New(ft.Type)
			dict, _ := fo.Dict()
			err := unmarshalDict(valp, dict)
			if err != nil {
				break
			}
			fv.Set(valp.Elem())
		}
	}
	return nil
}

func marshalValue(w io.Writer, v reflect.Value) int {
	len := 0
	switch v.Kind() {
	case reflect.String:
		len += EncodeString(w, v.String())
	case reflect.Int:
		len += EncodeInt(w, int(v.Int()))
	case reflect.Slice:
		len += marshalList(w, v)
	case reflect.Struct:
		len += marshalDict(w, v)
	}
	return len
}

func marshalList(w io.Writer, vl reflect.Value) int {
	len := 2
	w.Write([]byte{'l'})
	for i := 0; i < vl.Len(); i++ {
		ev := vl.Index(i)
		len += marshalValue(w, ev)
	}
	w.Write([]byte{'e'})
	return len
}

func marshalDict(w io.Writer, vd reflect.Value) int {
	len := 2
	w.Write([]byte{'d'})
	for i:=0; i < vd.NumField(); i++ {
		fv := vd.Field(i)
		ft := vd.Type().Field(i)
		len += EncodeString(w, strings.ToLower(ft.Name))
		len += marshalValue(w, fv)
	}
	w.Write([]byte{'e'})
	return len
}

func Marshal(w io.Writer, s interface{}) int {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return marshalValue(w, v)
}
