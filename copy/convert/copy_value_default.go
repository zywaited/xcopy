package convert

import (
	"reflect"

	"github.com/zywaited/xcopy/internal"
)

func init() {
	acv[reflect.Invalid] =
		newDefaultSearchValue(
			aliseSearchValue(newDefaultSearchValue(nil)),
		)
}

func newDefaultSearchValue(next ActualValue) ActualValue {
	return defaultMethodSearchValue( // field method
		defaultSearchValue( // Get{field} method
			next,
		),
	)
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
		mn := acFieldMethodNamePrefix + internal.ToCame(data.GetSf())
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
		mn := internal.ToCame(data.GetSf())
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
