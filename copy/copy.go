package copy

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/zywaited/xcopy/copy/convert"
	"github.com/zywaited/xcopy/copy/option"

	"github.com/pkg/errors"
)

type XCopier interface {
	Copy(dest, source interface{}) (err error)
}

// Deprecated
type XCopy interface {
	// CopySF 结构体赋值方法【这里抽离是为了兼容已有逻辑，可以使用CopyF替代】
	CopySF(dest, source interface{}) (err error)
	// CopyF 通用的赋值方法
	CopyF(dest, source interface{}) (err error)
}

type kindCopier interface {
	check(dv, sv reflect.Value) error
	copy(c *xCopy, dv, sv reflect.Value) (err error)
	recursion(c *xCopy, data, nd *convert.Info) bool
}

var kindCopiers = map[reflect.Kind]kindCopier{}

// 字段赋值
// 结构体赋值时对应字段可指定源的字段(`copy:"source's field"`)
type xCopy struct {
	// 强制转换(慎用)
	// 多级指针只有层级一致或者dv和sv最后都不再是指针才不会有问题
	// Recursion时, 多级指针层级一致或者sv最后不再是指针才不会有问题
	convert bool

	// 出错是否继续下一个
	next bool

	// 是否递归(结构体、数组)
	// 不过这样也会让赋值都变成值赋值而不是遇到指针是引用，也就是部分指针会重新申请内存，类似于深度拷贝
	// 并且只有在Convert为true时生效
	// 递归有性能损耗，并且如果赋值中有循环依赖可能导致死循环，与json解析一样，慎用
	recursion bool

	// 当copy赋值字段为空时是否需要读取JSON TAG中的字段
	jsonTag bool

	xcm         convert.XConverters
	acv         convert.ActualValuer
	cp          *sync.Pool
	kindCopiers map[reflect.Kind]kindCopier
}

const (
	OriginCopyField     = "origin"
	FuncCopyFieldPrefix = "func"
)

// NewCopy 初始化默认
func NewCopy(opts ...option.Option) *xCopy {
	c := option.Config{
		Convert:   true,
		Next:      true,
		Recursion: true,
		JsonTag:   true,
		Xcm:       convert.AcDefaultXConverter(),
		Acv:       convert.AcDefaultActualValuer(),
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &xCopy{
		convert:   c.Convert,
		next:      c.Next,
		recursion: c.Recursion,
		jsonTag:   c.JsonTag,
		xcm:       c.Xcm,
		acv:       c.Acv,
		cp: &sync.Pool{New: func() interface{} {
			return &convert.Info{}
		}},
		kindCopiers: kindCopiers,
	}
}

// Clone 克隆
func (c *xCopy) Clone() *xCopy {
	return &xCopy{
		convert:     c.convert,
		next:        c.next,
		recursion:   c.recursion,
		jsonTag:     c.jsonTag,
		xcm:         c.xcm,
		acv:         c.acv,
		cp:          c.cp,
		kindCopiers: c.kindCopiers,
	}
}

// SetConvert 是否强转
func (c *xCopy) SetConvert(convert bool) *xCopy {
	cp := c.Clone()
	cp.convert = convert
	return cp
}

// SetNext 出错是否继续赋值下一个字段
func (c *xCopy) SetNext(next bool) *xCopy {
	cp := c.Clone()
	cp.next = next
	return cp
}

// SetRecursion 是否递归（依赖强转）
func (c *xCopy) SetRecursion(recursion bool) *xCopy {
	cp := c.Clone()
	cp.recursion = recursion
	return cp
}

// SetJSONTag 是否读取JSON TAG
func (c *xCopy) SetJSONTag(jsonTag bool) *xCopy {
	cp := c.Clone()
	cp.jsonTag = jsonTag
	return cp
}

// 抽离函数，过滤不合法的字段
func (c *xCopy) setSf(data *convert.Info) (err error) {
	// 强制转换可能会出现异常
	defer func() {
		if pe := recover(); pe != nil {
			err = fmt.Errorf("赋值失败: [%#v]", pe)
		}
	}()

	// 检车字段是否有效
	if data.GetSf() == "" {
		data.SetSf(data.GetDf())
	}
	if data.GetDf() == "" && data.GetSf() == "" || data.GetSf() == "-" {
		return
	}

	dFieldValue := data.GetDv().FieldByName(data.GetDf())
	// 判断是否结构体字段存在
	if !dFieldValue.IsValid() || !dFieldValue.CanSet() {
		return
	}
	// 源字段检查
	stKind := data.GetSv().Kind()
	ach := c.acv.AC(stKind)
	err = ach(data)
	if err != nil {
		return err
	}
	sFieldValue := data.GetSv()
	if !sFieldValue.IsValid() {
		err = errors.New("赋值类型必须是struct、map或者array(slice); 源字段不存在")
		return
	}
	data.SetDv(dFieldValue)
	data.SetSv(sFieldValue)
	c.value(data)
	return
}

// NOTE 复杂逻辑不符合预期: 递归结构体、多级指针
// 抽离函数，真正的赋值操作
func (c *xCopy) value(data *convert.Info) bool {
	dt := data.GetDv().Type()
	st := data.GetSv().Type()
	// nil返回
	if (st.Kind() == reflect.Ptr || st.Kind() == reflect.Interface) && data.GetSv().IsNil() {
		return false
	}

	// 判断类型是否一样或者被赋值是interface
	if dt == st || dt.Kind() == reflect.Interface {
		data.GetDv().Set(data.GetSv())
		return true
	}

	// 如果赋值是interface
	if st.Kind() == reflect.Interface {
		data.SetSv(reflect.ValueOf(data.GetSv().Interface()))
		return c.value(data)
	}

	// 是否强制转换
	if !c.convert {
		return false
	}

	// 找到未初始化的指针
	for dt.Kind() == reflect.Ptr && !data.GetDv().IsNil() {
		if st.Kind() == reflect.Ptr {
			if data.GetSv().IsNil() {
				return false
			}
			data.SetSv(data.GetSv().Elem())
			st = data.GetSv().Type()
		}
		data.SetDv(data.GetDv().Elem())
		dt = data.GetDv().Type()
	}

	if !c.recursion {
		return c.notRecursion(data)
	}
	// 指针赋值(尽可能重新申请内存拷贝)
	nd, next := c.recursionPointer(data)
	if nd == nil {
		return false
	}
	defer c.cp.Put(nd)
	if !next {
		// 那么可以直接赋值
		if data.GetDv() != nd.GetSv() {
			data.GetDv().Set(nd.GetSv())
		}
		return false
	}
	// 指针过后类型还是不同的话，这里做个判断是为了防止死循环
	dt = nd.GetDv().Type()
	st = data.GetSv().Type()
	if st != dt {
		kc := c.kindCopiers[dt.Kind()]
		if kc == nil || kc.check(nd.GetDv(), data.GetSv()) != nil {
			kc = c.kindCopiers[reflect.Invalid]
		}
		return kc.recursion(c, data, nd)
	}
	return true
}

// 不递归进行指针赋值
func (c *xCopy) notRecursion(data *convert.Info) bool {
	dt := data.GetDv().Type()
	st := data.GetSv().Type()
	// 赋值指针(不申请内存拷贝)
	if dt.Kind() == reflect.Ptr {
		// 转换赋值为指针
		if st.Kind() != reflect.Ptr {
			sv := reflect.New(st)
			sv.Elem().Set(data.GetSv())
			data.SetSv(sv)
		}
		if data.GetSv().IsNil() {
			return false
		}
		// 赋值
		data.GetDv().Set(reflect.NewAt(dt.Elem(), unsafe.Pointer(data.GetSv().Pointer())))
		return true
	}
	// 多级指针
	for st.Kind() == reflect.Ptr {
		if data.GetSv().IsNil() {
			return false
		}
		data.SetSv(data.GetSv().Elem())
		st = data.GetSv().Type()
	}
	if st == dt {
		data.GetDv().Set(data.GetSv())
		return true
	}
	return c.kindCopiers[reflect.Invalid].recursion(c, data, data)
}

// 递归指针赋值
func (c *xCopy) recursionPointer(data *convert.Info) (*convert.Info, bool) {
	dt := data.GetDv().Type()
	st := data.GetSv().Type()
	malloc := false

	nd := c.cp.Get().(*convert.Info)
	// 下一步的赋值全是dv
	nd.SetDv(data.GetDv())
	nd.SetSv(data.GetDv())
	nd.SetDf(data.GetDf())
	nd.SetSf(data.GetSf())
	nd.SetOfn(data.IsOfn())
	nd.SetOsf(data.IsOsf())

	for {
		if st == dt {
			nd.GetDv().Set(data.GetSv())
			return nd, false
		}
		if dt.Kind() != reflect.Ptr && st.Kind() != reflect.Ptr {
			break
		}
		if dt.Kind() == reflect.Ptr {
			if !malloc {
				// 这里要独立出数据(因为后面可能不需要赋值)
				malloc = true
				nd.SetDv(reflect.New(dt.Elem()))
				nd.SetSv(nd.GetDv())
			} else {
				nd.GetDv().Set(reflect.New(dt.Elem()))
			}
			nd.SetDv(reflect.Indirect(nd.GetDv()))
			dt = nd.GetDv().Type()
		}
		if st.Kind() == reflect.Ptr {
			if data.GetSv().IsNil() {
				// 没办法赋值
				return nil, false
			}
			data.SetSv(reflect.Indirect(data.GetSv()))
			st = data.GetSv().Type()
		}
	}
	return nd, true
}

// 递归赋值数组
func (c *xCopy) recursionSlice(data, nd *convert.Info) bool {
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

// Deprecated
func (c *xCopy) CopySF(dest, source interface{}) (err error) {
	return c.Copy(dest, source)
}

// Deprecated
func (c *xCopy) CopyF(dest, source interface{}) (err error) {
	return c.Copy(dest, source)
}

// Copy CopyF
func (c *xCopy) Copy(dest, source interface{}) error {
	if source == nil {
		return errors.New("赋值体不存在")
	}
	// 校验类型
	dv := reflect.ValueOf(dest)
	if dv.Type().Kind() != reflect.Ptr {
		return errors.New("被赋值的单体必须是指针类型")
	}

	dv = dv.Elem()
	sv := reflect.Indirect(reflect.ValueOf(source))
	kc := c.kindCopiers[dv.Type().Kind()]
	if kc == nil {
		kc = c.kindCopiers[reflect.Invalid]
	}
	return kc.copy(c, dv, sv)
}
