package copy

import (
	"reflect"
	"strings"

	"github.com/zywaited/xcopy/copy/convert"
)

func init() {
	kindCopiers[reflect.Invalid] = newDefaultCopier()
}

func newDefaultCopier() *defaultCopier {
	return &defaultCopier{}
}

type defaultCopier struct {
}

func (dc *defaultCopier) check(_, _ reflect.Value) error {
	return nil
}

func (dc *defaultCopier) copy(c *xCopy, dv, sv reflect.Value) error {
	dv = reflect.Indirect(dv)
	sv = reflect.Indirect(sv)
	// 重置数据
	data := c.cp.Get().(*convert.Info)
	defer c.cp.Put(data)
	data.SetDf("")
	data.SetSf("")
	data.SetOsf(false)
	data.SetOfn(false)
	data.SetDv(dv)
	data.SetSv(sv)
	c.value(data)
	return nil
}

// 强制转换: 调用该方法的时候已经不再是指针
func (dc *defaultCopier) recursion(c *xCopy, data, nd *convert.Info) bool {
	dt := nd.GetDv().Type()
	sv := data.GetSv()
	nsv := nd.GetSv()
	vh := c.xcm.AC(strings.TrimLeft(dt.PkgPath()+"."+dt.Name(), "."))
	if vh != nil {
		nd.SetSv(data.GetSv())
		sv = vh.Convert(nd)
	}
	nd.GetDv().Set(sv.Convert(dt))
	if data.GetDv() != nsv {
		data.GetDv().Set(nsv)
	}
	return true
}
