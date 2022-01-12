package convert

import (
	"errors"
	"reflect"
)

type (
	ActualValue func(*Info) error

	ActualValueMs map[reflect.Kind]ActualValue
)

type ActualValuer interface {
	AC(reflect.Kind) ActualValue
}

func (ms ActualValueMs) AC(kind reflect.Kind) ActualValue {
	if ms[kind] != nil {
		return ms[kind]
	}
	return ms[reflect.Invalid]
}

func (ms ActualValueMs) Clone() ActualValueMs {
	c := ActualValueMs{}
	for kind, actualValuer := range ms {
		c[kind] = actualValuer
	}
	return c
}

var (
	acv ActualValueMs = map[reflect.Kind]ActualValue{}

	invalidValue  = errors.New("目标数据无效")
	invalidMethod = errors.New("目标函数无效")
)

func AcDefaultActualValuer() ActualValuer {
	return acv
}
