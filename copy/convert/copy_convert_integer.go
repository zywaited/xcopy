package convert

import (
	"reflect"
	"strconv"
	"time"
)

func init() {
	// int uint
	ic := NewIntXConverter(dc)
	for _, in := range []string{"", "8", "16", "32", "64"} {
		xcms.Register("int"+in, ic)
		xcms.Register("uint"+in, ic)
	}
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
	if kind == reflect.Struct && st.PkgPath()+"."+st.Name() == timeConverterName {
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
