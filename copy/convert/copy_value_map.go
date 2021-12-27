package convert

import (
	"reflect"

	"github.com/zywaited/xcopy/utils"
)

func init() {
	acv[reflect.Map] =
		methodSearchValue(
			newMapSearchValue(
				aliseSearchValue(newMapSearchValue(nil)),
			),
		)
}

func newMapSearchValue(next ActualValue) ActualValue {
	return mapSearchValue( // map filed
		newDefaultSearchValue(next), // default
	)
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
		for i := 0; i < 4; i++ {
			sf := data.GetSf()
			switch i {
			case 1:
				sf = utils.ToSnake(sf)
			case 2:
				sf = utils.ToLcFirst(utils.ToCame(sf))
			case 3:
				sf = utils.ToCame(sf)
			}
			v := data.GetSv().MapIndex(reflect.ValueOf(sf))
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
