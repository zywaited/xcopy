package convert

import (
	"reflect"
	"strconv"
	"time"
)

func init() {
	xcms.Register("string", NewStringXConverter(dc))
}

type StringXConverter struct {
	next XConverter
}

func NewStringXConverter(next XConverter) *StringXConverter {
	return &StringXConverter{next: next}
}

func (sc *StringXConverter) Convert(data *Info) reflect.Value {
	sv := data.GetSv()
	if !sv.IsValid() {
		return sv
	}
	st := sv.Type()
	kind := st.Kind()
	switch kind {
	case reflect.Struct:
		if st.PkgPath() == "time" && st.Name() == "Time" {
			return reflect.ValueOf(sv.Interface().(time.Time).Format("2006-01-02 15:04:05"))
		}
	case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(strconv.FormatInt(sv.Int(), 10))
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(strconv.FormatUint(sv.Uint(), 10))
	}
	if sc.next != nil {
		sv = sc.next.Convert(data)
	}
	if sv != data.GetSv() {
		return sv
	}
	for _, name := range []string{"String", "ToString"} {
		mv := data.GetSv().MethodByName(name)
		if !mv.IsValid() {
			if !data.GetSv().CanAddr() {
				continue
			}
			mv = data.GetSv().Addr().MethodByName(name)
			if !mv.IsValid() {
				continue
			}
		}
		mt := mv.Type()
		if mt.NumIn() > 0 || mt.NumOut() == 0 {
			continue
		}
		rsv := mv.Call(nil)[0]
		// 因为是string函数，返回值按规范必须是string
		if !rsv.IsValid() || rsv.Kind() != reflect.String {
			continue
		}
		return rsv
	}
	return sv
}
