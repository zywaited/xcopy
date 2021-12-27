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

func (sc *StringXConverter) Convert(data *Info) (sv reflect.Value) {
	sv = data.GetSv()
	if !data.GetSv().IsValid() {
		return
	}
	rk := true
	defer func() {
		if !rk && sc.next != nil {
			sv = sc.next.Convert(data)
		}
	}()
	st := data.GetSv().Type()
	kind := st.Kind()
	switch kind {
	case reflect.Struct:
		if st.PkgPath() == "time" && st.Name() == "Time" {
			sv = reflect.ValueOf(sv.Interface().(time.Time).Format("2006-01-02 15:04:05"))
			return
		}
	case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
		sv = reflect.ValueOf(strconv.FormatInt(sv.Int(), 10))
		return
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		sv = reflect.ValueOf(strconv.FormatUint(sv.Uint(), 10))
		return
	}
	rk = false
	mv := data.GetSv().MethodByName("String")
	if !mv.IsValid() {
		if !data.GetSv().CanAddr() {
			return
		}
		mv = data.GetSv().Addr().MethodByName("String")
		if !mv.IsValid() {
			return
		}
	}
	mt := mv.Type()
	if mt.NumIn() > 0 || mt.NumOut() == 0 {
		return
	}
	rsv := mv.Call(nil)[0]
	// 因为是string函数，返回值按规范必须是string
	if !rsv.IsValid() || rsv.Kind() != reflect.String {
		return
	}
	rk = true
	sv = rsv
	return
}
