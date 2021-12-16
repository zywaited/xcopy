package xcopy

import (
	"errors"
	"reflect"
	"strconv"
)

type (
	ActualValue func(*convertInfo) error

	ActualValueMs map[reflect.Kind]ActualValue
)

var (
	acv ActualValueMs

	invalidValue = errors.New("目标数据无效")
)

func init() {
	acv = make(map[reflect.Kind]ActualValue)
	acv[reflect.Array] = arraySearchValue(defaultSearchValue(nil))
	acv[reflect.Slice] = arraySearchValue(defaultSearchValue(nil))
	acv[reflect.Map] = mapSearchValue(defaultSearchValue(nil))
	acv[reflect.Struct] = stSearchValue(defaultSearchValue(nil))
}

// 默认读取
func defaultSearchValue(next ActualValue) ActualValue {
	return func(data *convertInfo) (err error) {
		if !data.sv.IsValid() {
			err = invalidValue
			return
		}
		defer func() {
			if err != nil && next != nil {
				err = next(data)
			}
		}()
		mv := data.sv.MethodByName(acFieldMethodNamePrefix + ToCame(data.sf))
		if !mv.IsValid() {
			err = invalidValue
			return
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
		data.sv = v
		return
	}
}

func arraySearchValue(next ActualValue) ActualValue {
	return func(data *convertInfo) (err error) {
		if !data.sv.IsValid() {
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
		d, err = strconv.Atoi(data.sf)
		if err != nil {
			return
		}
		v := data.sv.Index(d)
		if !v.IsValid() {
			err = invalidValue
			return
		}
		data.sv = v
		return
	}
}

func mapSearchValue(next ActualValue) ActualValue {
	return func(data *convertInfo) (err error) {
		if !data.sv.IsValid() {
			err = invalidValue
			return
		}
		defer func() {
			if err != nil && next != nil {
				err = next(data)
			}
		}()
		sf := data.sf
		if !data.osf {
			sf = ToSnake(sf)
		}
		v := data.sv.MapIndex(reflect.ValueOf(sf))
		if !v.IsValid() {
			err = invalidValue
			return
		}
		data.sv = v
		return
	}
}

func stSearchValue(next ActualValue) ActualValue {
	return func(data *convertInfo) (err error) {
		if !data.sv.IsValid() {
			err = invalidValue
			return
		}
		defer func() {
			if err != nil && next != nil {
				err = next(data)
			}
		}()
		sf := data.sf
		if !data.osf {
			sf = ToCame(sf)
		}
		v := data.sv.FieldByName(sf)
		if !v.IsValid() {
			err = invalidValue
			return
		}
		data.sv = v
		return
	}
}
