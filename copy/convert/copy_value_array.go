package convert

import (
	"reflect"
	"strconv"
)

func init() {
	action :=
		methodSearchValue(
			newArraySearchValue(
				toArraySearchValue(nil),
			),
		)
	acv[reflect.Array] = action
	acv[reflect.Slice] = action
}

func newArraySearchValue(next ActualValue) ActualValue {
	return arraySearchValue( // array[slice] index
		newDefaultSearchValue(next), // default
	)
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

func toArraySearchValue(next ActualValue) ActualValue {
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
		mn := "ToArray"
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
