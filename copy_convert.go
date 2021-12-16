package xcopy

import (
	"reflect"
	"strconv"
	"time"
)

const acFieldMethodNamePrefix = "Get"

type (
	convertInfo struct {
		df  string
		sf  string
		osf bool
		dv  reflect.Value
		sv  reflect.Value
	}

	XConverter interface {
		Convert(*convertInfo) reflect.Value
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

// 默认都要走的处理器
var dc = NewDefaultXConverter()

// 默认的转换器
type defaultXConverter struct {
}

func NewDefaultXConverter() *defaultXConverter {
	return &defaultXConverter{}
}

func (dc *defaultXConverter) Convert(data *convertInfo) reflect.Value {
	if data.sf == "" || !data.sv.IsValid() {
		return data.sv
	}
	mn := acFieldMethodNamePrefix + ToCame(data.sf)
	mv := data.sv.MethodByName(mn)
	if !mv.IsValid() {
		if !data.sv.CanAddr() {
			return data.sv
		}
		mv = data.sv.Addr().MethodByName(mn)
		if !mv.IsValid() {
			return data.sv
		}
	}
	mt := mv.Type()
	if mt.NumIn() > 0 || mt.NumOut() == 0 {
		return data.sv
	}
	sv := mv.Call(nil)[0]
	if !sv.IsValid() {
		return data.sv
	}
	return sv
}

type timeXConverter struct {
	next XConverter
}

func NewTimeXConverter(next XConverter) *timeXConverter {
	return &timeXConverter{next: next}
}

func (tc *timeXConverter) Convert(data *convertInfo) reflect.Value {
	if !data.sv.IsValid() {
		return data.sv
	}
	sv := data.sv
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

func (ic *IntXConverter) Convert(data *convertInfo) (sv reflect.Value) {
	sv = data.sv
	if !data.sv.IsValid() {
		return
	}
	rk := true
	defer func() {
		if !rk && ic.next != nil {
			sv = ic.next.Convert(data)
		}
	}()
	st := data.sv.Type()
	kind := st.Kind()
	if kind == reflect.Struct && st.PkgPath() == "time" && st.Name() == "Time" {
		sv = reflect.ValueOf(data.sv.Interface().(time.Time).Unix())
		return
	}
	if kind != reflect.String {
		rk = false
		return
	}
	switch data.dv.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		d, err := strconv.Atoi(data.sv.String())
		if err != nil {
			rk = false
			return
		}
		sv = reflect.ValueOf(d)
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		d, err := strconv.ParseUint(data.sv.String(), 10, 0)
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

func (sc *StringXConverter) Convert(data *convertInfo) (sv reflect.Value) {
	sv = data.sv
	if !data.sv.IsValid() {
		return
	}
	rk := true
	defer func() {
		if !rk && sc.next != nil {
			sv = sc.next.Convert(data)
		}
	}()
	st := data.sv.Type()
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
	mv := data.sv.MethodByName("String")
	if !mv.IsValid() {
		if !data.sv.CanAddr() {
			return
		}
		mv = data.sv.Addr().MethodByName("String")
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
