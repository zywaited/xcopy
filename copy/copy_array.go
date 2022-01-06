package copy

import (
	"reflect"

	"github.com/zywaited/xcopy/copy/convert"

	"github.com/pkg/errors"
)

func init() {
	kindCopiers[reflect.Array] = newArrayCopier()
	kindCopiers[reflect.Slice] = newArrayCopier()
}

func newArrayCopier() *arrayCopier {
	return &arrayCopier{}
}

type arrayCopier struct {
}

func (ac *arrayCopier) check(dv, sv reflect.Value) error {
	dk := reflect.Indirect(dv).Kind()
	if dk != reflect.Array && dk != reflect.Slice {
		return errors.New("被赋值的不是数组")
	}

	sk := reflect.Indirect(sv).Kind()
	if sk != reflect.Array && sk != reflect.Slice {
		return errors.New("赋值体的不是数组")
	}

	return nil
}

func (ac *arrayCopier) copy(c *xCopy, dv, sv reflect.Value) (err error) {
	sl := sv.Len()
	if sl == 0 {
		return
	}
	dt := dv.Type()
	if dt.Kind() == reflect.Array {
		dl := dv.Len()
		if dl < sl {
			sl = dl
		}
	} else {
		dv.Set(reflect.MakeSlice(dt, sl, sl))
	}
	data := c.cp.Get().(*convert.Info)
	defer c.cp.Put(data)
	data.SetDf("")
	data.SetSf("")
	for i := 0; i < sl; i++ {
		data.SetDv(dv.Index(i))
		data.SetSv(sv.Index(i))
		if !c.value(data) {
			continue
		}
	}
	return
}

func (ac *arrayCopier) recursion(c *xCopy, data, nd *convert.Info) bool {
	sl := data.GetSv().Len()
	if sl == 0 {
		return false
	}
	tdt := nd.GetDv().Type()
	if tdt.Kind() == reflect.Array {
		dl := nd.GetDv().Len()
		if dl < sl {
			sl = dl
		}
	} else {
		data.GetDv().Set(reflect.MakeSlice(tdt, sl, sl))
	}
	ok := sl == 0 // 等于0算成功
	ndv := nd.GetDv()
	sv := nd.GetSv()
	for i := 0; i < sl; i++ {
		nd.SetDv(ndv.Index(i))
		nd.SetSv(data.GetSv().Index(i))
		if !c.value(nd) {
			continue
		}
		ok = true
	}
	if ok && data.GetDv() != sv {
		data.GetDv().Set(sv)
	}
	return ok
}
