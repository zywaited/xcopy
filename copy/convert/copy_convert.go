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

func AcDefaultXConverter() XConverters {
	return xcms
}
