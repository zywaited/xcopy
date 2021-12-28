package convert

import (
	"reflect"

	"github.com/zywaited/xcopy/internal"
)

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
	mn := acFieldMethodNamePrefix + internal.ToCame(data.GetSf())
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
	mn := internal.ToCame(data.GetSf())
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
