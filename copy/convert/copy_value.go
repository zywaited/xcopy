package convert

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/zywaited/xcopy/utils"
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
	return acvDC
}

var (
	acv ActualValueMs

	invalidValue  = errors.New("目标数据无效")
	invalidMethod = errors.New("目标函数无效")
)

func init() {
	acv = make(map[reflect.Kind]ActualValue)
	acv[reflect.Array] = methodSearchValue(arraySearchValue(acvDC))
	acv[reflect.Slice] = methodSearchValue(arraySearchValue(acvDC))
	acv[reflect.Map] = methodSearchValue(mapSearchValue(acvDC))
	acv[reflect.Struct] = methodSearchValue(stSearchValue(acvDC))
}

var acvDC = defaultMethodSearchValue(defaultSearchValue(nil))

func AcDefaultActualValuer() ActualValuer {
	return acv
}

// 默认读取
func defaultSearchValue(next ActualValue) ActualValue {
	return func(data *Info) (err error) {
		if !data.GetSv().IsValid() {
			err = invalidValue
			return
		}
		defer func() {
			if err != nil && next != nil {
				err = next(data)
			}
		}()
		mn := acFieldMethodNamePrefix + utils.ToCame(data.GetSf())
		mv := data.GetSv().MethodByName(mn)
		if !mv.IsValid() {
			if !data.GetSv().CanAddr() {
				return invalidValue
			}
			mv = data.GetSv().Addr().MethodByName(mn)
			if !mv.IsValid() {
				return invalidValue
			}
		}
		mt := mv.Type()
		if mt.NumIn() > 0 || mt.NumOut() == 0 {
			err = invalidValue
			return
		}
		v := mv.Call(nil)[0]
		if !v.IsValid() {
			err = invalidValue
			return
		}
		data.SetSv(v)
		return
	}
}

func defaultMethodSearchValue(next ActualValue) ActualValue {
	return func(data *Info) (err error) {
		if !data.GetSv().IsValid() {
			err = invalidValue
			return
		}
		defer func() {
			if err != nil && next != nil {
				err = next(data)
			}
		}()
		mn := utils.ToCame(data.GetSf())
		mv := data.GetSv().MethodByName(mn)
		if !mv.IsValid() {
			if !data.GetSv().CanAddr() {
				return invalidValue
			}
			mv = data.GetSv().Addr().MethodByName(mn)
			if !mv.IsValid() {
				return invalidValue
			}
		}
		mt := mv.Type()
		if mt.NumIn() > 0 || mt.NumOut() == 0 {
			err = invalidValue
			return
		}
		v := mv.Call(nil)[0]
		if !v.IsValid() {
			err = invalidValue
			return
		}
		data.SetSv(v)
		return
	}
}

func arraySearchValue(next ActualValue) ActualValue {
	return func(data *Info) (err error) {
		if !data.GetSv().IsValid() {
			err = invalidValue
			return
		}
		defer func() {
			if err != nil && next != nil {
				err = next(data)
			}
		}()
		// 解析sf为数字
		d := 0
		d, err = strconv.Atoi(data.GetSf())
		if err != nil {
			return
		}
		v := data.GetSv().Index(d)
		if !v.IsValid() {
			err = invalidValue
			return
		}
		data.SetSv(v)
		return
	}
}

func mapSearchValue(next ActualValue) ActualValue {
	return func(data *Info) (err error) {
		if !data.GetSv().IsValid() {
			err = invalidValue
			return
		}
		defer func() {
			if err != nil && next != nil {
				err = next(data)
			}
		}()
		sf := data.GetSf()
		if !data.IsOsf() {
			sf = utils.ToSnake(sf)
		}
		v := data.GetSv().MapIndex(reflect.ValueOf(sf))
		if !v.IsValid() {
			err = invalidValue
			return
		}
		data.SetSv(v)
		return
	}
}

func stSearchValue(next ActualValue) ActualValue {
	return func(data *Info) (err error) {
		if !data.GetSv().IsValid() {
			err = invalidValue
			return
		}
		defer func() {
			if err != nil && next != nil {
				err = next(data)
			}
		}()
		sf := data.GetSf()
		if !data.IsOsf() {
			sf = utils.ToCame(sf)
		}
		v := data.GetSv().FieldByName(sf)
		if !v.IsValid() {
			err = invalidValue
			return
		}
		data.SetSv(v)
		return
	}
}

func methodSearchValue(next ActualValue) ActualValue {
	return func(data *Info) error {
		if !data.IsOfn() && next != nil {
			return next(data)
		}
		if !data.IsOfn() && next == nil {
			return invalidValue
		}
		method := data.GetSf()
		if !data.IsOsf() {
			method = utils.ToCame(method)
		}
		mv := data.GetSv().MethodByName(method)
		if !mv.IsValid() {
			if !data.GetSv().CanAddr() {
				return invalidMethod
			}
			mv = data.GetSv().Addr().MethodByName(method)
			if !mv.IsValid() {
				return invalidMethod
			}
		}
		mt := mv.Type()
		if mt.NumIn() > 0 || mt.NumOut() == 0 {
			return invalidMethod
		}
		// 这里要重置字段名称
		data.SetSf(data.GetDf())
		data.SetSv(mv.Call(nil)[0])
		return nil
	}
}
