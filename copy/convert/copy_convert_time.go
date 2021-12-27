package convert

import (
	"reflect"
	"time"
)

func init() {
	xcms.Register("time.Time", NewTimeXConverter(dc))
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
	}
	return sv
}
