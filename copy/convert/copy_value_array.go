package convert

import (
	"reflect"
	"strconv"
)

func init() {
	action :=
		methodSearchValue(
			newArraySearchValue(nil),
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
