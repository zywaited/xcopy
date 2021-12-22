package convert

import (
	"reflect"
	"strconv"
	"time"

	"github.com/zywaited/xcopy/utils"
)

const acFieldMethodNamePrefix = "Get"

type (
	XConverter interface {
		Convert(*Info) reflect.Value
	}

	XConverters interface {
		AC(string) XConverter
	}

	NConverterM map[string]XConverter
)

var xcms = make(NConverterM)

func (xcm NConverterM) AC(name string) XConverter {
	if xcm[name] == nil {
		return dc
	}
	return xcm[name]
}

func (xcm NConverterM) Register(name string, xc XConverter) {
	xcm[name] = xc
}

func init() {
	xcms.Register("time.Time", NewTimeXConverter(dc))
	xcms.Register("string", NewStringXConverter(dc))
	// int uint
	ic := NewIntXConverter(dc)
	for _, in := range []string{"", "8", "16", "32", "64"} {
		xcms.Register("int"+in, ic)
		xcms.Register("uint"+in, ic)
	}
}

func AcDefaultXConverter() XConverters {
	return xcms
}

// 默认都要走的处理器
var dc = NewDefaultMethodXConverter(NewDefaultXConverter())

// 默认的转换器
type defaultXConverter struct {
}

func NewDefaultXConverter() *defaultXConverter {
	return &defaultXConverter{}
}

func (dc *defaultXConverter) Convert(data *Info) reflect.Value {
	if data.GetSf() == "" || !data.GetSv().IsValid() {
		return data.GetSv()
	}
	mn := acFieldMethodNamePrefix + utils.ToCame(data.GetSf())
	mv := data.GetSv().MethodByName(mn)
	if !mv.IsValid() {
		if !data.GetSv().CanAddr() {
			return data.GetSv()
		}
		mv = data.GetSv().Addr().MethodByName(mn)
		if !mv.IsValid() {
			return data.GetSv()
		}
	}
	mt := mv.Type()
	if mt.NumIn() > 0 || mt.NumOut() == 0 {
		return data.GetSv()
	}
	sv := mv.Call(nil)[0]
	if !sv.IsValid() {
		return data.GetSv()
	}
	return sv
}

type defaultMethodXConverter struct {
	next XConverter
}

func NewDefaultMethodXConverter(next XConverter) *defaultMethodXConverter {
	return &defaultMethodXConverter{next: next}
}

func (dm *defaultMethodXConverter) Convert(data *Info) (sv reflect.Value) {
	sv = data.GetSv()
	if data.GetSf() == "" || !sv.IsValid() {
		return
	}
	rk := true
	defer func() {
		if !rk && dm.next != nil {
			sv = dm.next.Convert(data)
		}
	}()
	mn := utils.ToCame(data.GetSf())
	mv := data.GetSv().MethodByName(mn)
	if !mv.IsValid() {
		if !data.GetSv().CanAddr() {
			rk = false
			return
		}
		mv = data.GetSv().Addr().MethodByName(mn)
		if !mv.IsValid() {
			rk = false
			return
		}
	}
	mt := mv.Type()
	if mt.NumIn() > 0 || mt.NumOut() == 0 {
		rk = false
		return
	}
	sv = mv.Call(nil)[0]
	return
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

type IntXConverter struct {
	next XConverter
}

func NewIntXConverter(next XConverter) *IntXConverter {
	return &IntXConverter{next: next}
}

func (ic *IntXConverter) Convert(data *Info) (sv reflect.Value) {
	sv = data.GetSv()
	if !sv.IsValid() {
		return
	}
	rk := true
	defer func() {
		if !rk && ic.next != nil {
			sv = ic.next.Convert(data)
		}
	}()
	st := data.GetSv().Type()
	kind := st.Kind()
	if kind == reflect.Struct && st.PkgPath() == "time" && st.Name() == "Time" {
		sv = reflect.ValueOf(data.GetSv().Interface().(time.Time).Unix())
		return
	}
	if kind != reflect.String {
		rk = false
		return
	}
	switch data.GetDv().Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		d, err := strconv.Atoi(data.GetSv().String())
		if err != nil {
			rk = false
			return
		}
		sv = reflect.ValueOf(d)
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		d, err := strconv.ParseUint(data.GetSv().String(), 10, 0)
		if err != nil {
			rk = false
			return
		}
		sv = reflect.ValueOf(d)
	}
	return sv
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
