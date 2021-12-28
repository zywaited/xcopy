package convert

import (
	"reflect"
	"time"
)

const timeConverterName = "time.Time"

func init() {
	xcms.SkipCopier(timeConverterName)
	xcms.Register(timeConverterName, NewTimeXConverter(dc))
}

type timeXConverter struct {
	next XConverter
}

func NewTimeXConverter(next XConverter) *timeXConverter {
	return &timeXConverter{next: next}
}

func (tc *timeXConverter) Convert(data *Info) reflect.Value {
	if !data.GetSv().IsValid() {
		return data.GetSv()
	}
	sv := data.GetSv()
	switch sv.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		sv = reflect.ValueOf(time.Unix(sv.Int(), 0))
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		sv = reflect.ValueOf(time.Unix(int64(sv.Uint()), 0)) // 如果uint64超过int64，这里会抛异常
	case reflect.String:
		now, err := time.ParseInLocation("2006-01-02 15:04:05", sv.String(), time.Local)
		if err == nil {
			sv = reflect.ValueOf(now)
		}
	default:
		if tc.next != nil {
			sv = tc.next.Convert(data)
		}
		if sv != data.GetSv() {
			return sv
		}
		for _, name := range []string{"Time", "ToTime"} {
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
			if !rsv.IsValid() {
				continue
			}
			rst := rsv.Type()
			if rst.Kind() != reflect.Struct || rst.PkgPath()+"."+rst.Name() != timeConverterName {
				continue
			}
			return rsv
		}
	}
	return sv
}
