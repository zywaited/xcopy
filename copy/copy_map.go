package copy

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/zywaited/xcopy/copy/convert"

	"github.com/pkg/errors"
)

func init() {
	kindCopiers[reflect.Map] = newMapCopier()
}

var (
	invalidMapKeyType = errors.New("map key 不是字符串类型")
)

const omitempty = "omitempty"

func newMapCopier() *mapCopier {
	return &mapCopier{}

}

type mapCopier struct {
}

func (mc *mapCopier) check(dv, sv reflect.Value) error {
	dk := reflect.Indirect(dv).Kind()
	if dk != reflect.Map {
		return errors.New("被赋值的不是哈希")
	}

	sk := reflect.Indirect(sv).Kind()
	if sk != reflect.Map && sk != reflect.Struct {
		return errors.New("赋值体的不是结构体或者哈希")
	}
	return nil
}

func (mc *mapCopier) copy(c *xCopy, dv, sv reflect.Value) error {
	switch sv.Kind() {
	case reflect.Struct:
		sl := sv.NumField()
		if sl == 0 {
			return nil
		}
		dv.Set(reflect.MakeMapWithSize(dv.Type(), sl))
		return mc.structConvert(c, dv, sv)
	case reflect.Map:
		sl := sv.Len()
		if sl == 0 {
			return nil
		}
		dv.Set(reflect.MakeMapWithSize(dv.Type(), sl))
		return mc.mapConvert(c, dv, sv)
	}
	return nil
}

func (mc *mapCopier) recursion(c *xCopy, data, nd *convert.Info) bool {
	err := mc.copy(c, nd.GetDv(), data.GetSv())
	if err != nil {
		panic(errors.Wrap(err, "递归哈希赋值失败"))
	}
	if data.GetDv() != nd.GetSv() {
		data.GetDv().Set(nd.GetSv())
	}
	return true
}

// struct to map
func (mc *mapCopier) structConvert(c *xCopy, dv, sv reflect.Value) error {
	dt := dv.Type()
	// note: string
	if dt.Key().Kind() != reflect.String {
		return invalidMapKeyType
	}
	et := dt.Elem()
	st := sv.Type()
	num := st.NumField()
	var fst reflect.StructField
	data := c.cp.Get().(*convert.Info)
	defer c.cp.Put(data)
	data.SetDf("")
	data.SetSf("")
	data.SetOsf(false)
	data.SetOfn(false)
	for i := 0; i < num; i++ {
		fst = st.Field(i)
		name, opt := mc.parseTag(fst)
		if name == "-" {
			continue
		}
		esv := sv.FieldByName(fst.Name)
		if opt && mc.isEmptyValue(esv) {
			continue
		}
		edv := reflect.New(et)
		data.SetDv(edv)
		data.SetSv(esv)
		err := mc.set(c, data)
		if opt && err != nil {
			continue
		}
		if name == "" {
			name = fst.Name
		}
		if err != nil && !c.next {
			return errors.Wrapf(err, "赋值字段[%s]赋值失败", name)
		}
		if err == nil {
			dv.SetMapIndex(reflect.ValueOf(name), edv.Elem())
		}
	}
	return nil
}

func (mc *mapCopier) set(c *xCopy, data *convert.Info) (err error) {
	// 强制转换可能会出现异常
	defer func() {
		if pe := recover(); pe != nil {
			err = fmt.Errorf("赋值失败: [%s: %#v]", data.GetDf(), pe)
		}
	}()
	c.value(data)
	return
}

// struct tag
func (mc *mapCopier) parseTag(fst reflect.StructField) (string, bool) {
	sfs := strings.SplitN(strings.TrimSpace(fst.Tag.Get("json")), ",", 2)
	sf := strings.TrimSpace(sfs[0])
	oe := len(sfs) > 1 && strings.Contains(sfs[1], omitempty)
	return sf, oe
}

// json empty
func (mc *mapCopier) isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// map to map
func (mc *mapCopier) mapConvert(c *xCopy, dv, sv reflect.Value) error {
	dt := dv.Type()
	kt := dt.Key()
	et := dt.Elem()
	data := c.cp.Get().(*convert.Info)
	defer c.cp.Put(data)
	data.SetDf("")
	data.SetSf("")
	data.SetOsf(false)
	data.SetOfn(false)
	iter := sv.MapRange()
	for iter.Next() {
		dkv := reflect.New(kt)
		skv := iter.Key()
		data.SetDv(dkv)
		data.SetSv(skv)
		err := mc.set(c, data)
		if err != nil && !c.next {
			return errors.Wrapf(err, "赋值字段[%v]赋值失败", skv.Interface())
		}
		if err != nil {
			continue
		}
		dev := reflect.New(et)
		sev := iter.Value()
		data.SetDv(dev)
		data.SetSv(sev)
		err = mc.set(c, data)
		if err != nil && !c.next {
			return errors.Wrapf(err, "赋值字段[%v]赋值失败[Value]", skv.Interface())
		}
		if err == nil {
			dv.SetMapIndex(dkv.Elem(), dev.Elem())
		}
	}
	return nil
}
