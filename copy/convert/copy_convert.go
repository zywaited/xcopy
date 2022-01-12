package convert

import (
	"reflect"
)

const acFieldMethodNamePrefix = "Get"

type (
	XConverter interface {
		Convert(*Info) reflect.Value
	}

	XConverters interface {
		SC(string) bool
		AC(string) XConverter
	}

	XConvertersCloner interface {
		Clone() XConverters
	}

	XConvertersSetter interface {
		Register(name string, xc XConverter)
		SkipCopier(name string)
	}

	converterMS struct {
		m map[string]XConverter
		f map[string]bool
	}
)

func newConverterMS() *converterMS {
	return &converterMS{
		m: make(map[string]XConverter),
		f: make(map[string]bool),
	}
}

var xcms = newConverterMS()

func (xcm *converterMS) SC(name string) bool {
	return xcm.f[name]
}

func (xcm *converterMS) AC(name string) XConverter {
	if xcm.m[name] == nil {
		return dc
	}
	return xcm.m[name]
}

func (xcm *converterMS) Register(name string, xc XConverter) {
	xcm.m[name] = xc
}

func (xcm *converterMS) SkipCopier(name string) {
	xcm.f[name] = true
}

func (xcm *converterMS) Clone() XConverters {
	c := newConverterMS()
	for name, xConverter := range xcm.m {
		c.Register(name, xConverter)
	}
	for name, skip := range xcm.f {
		if skip {
			c.SkipCopier(name)
		}
	}
	return c
}

func AcDefaultXConverter() XConverters {
	return xcms
}
