package convert

import "strings"

func aliseSearchValue(next ActualValue) ActualValue {
	return func(data *Info) (err error) {
		if !data.GetSv().IsValid() {
			err = invalidValue
			return
		}
		sf := data.GetSf()
		osf := data.IsOsf()
		alise := strings.ToUpper(sf)
		if !commonInitialisms[alise] || next == nil {
			err = invalidValue
			return
		}
		// 执行完成后还原数据
		defer func() {
			data.SetSf(sf)
			data.SetOsf(osf)
		}()
		// 重置别名
		data.SetSf(alise)
		data.SetOsf(true)
		err = next(data)
		return
	}
}
