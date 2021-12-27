package convert

import (
	"reflect"

	"github.com/zywaited/xcopy/utils"
)

func init() {
	acv[reflect.Struct] =
		methodSearchValue( // ofn method
			newSTSearchValue( // filed
				aliseSearchValue(newSTSearchValue(nil)), // alise field
			),
		)
}

func newSTSearchValue(next ActualValue) ActualValue {
	return stSearchValue( // struct field
		newDefaultSearchValue(next), // default
	)
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
		for i := 0; i < 2; i++ {
			sf := data.GetSf()
			switch i {
			case 1:
				sf = utils.ToCame(sf)
			default:
				sf = utils.ToUcFirst(sf)
			}
			v := data.GetSv().FieldByName(sf)
			if v.IsValid() {
				data.SetSv(v)
				return
			}
			if data.IsOsf() {
				break
			}
		}
		err = invalidValue
		return
	}
}
