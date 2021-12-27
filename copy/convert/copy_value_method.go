package convert

import "github.com/zywaited/xcopy/internal"

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
			method = internal.ToCame(method)
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
