package xcopy

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unsafe"

	"github.com/pkg/errors"
)

type XCopy interface {
	// CopySF 结构体赋值方法【这里抽离是为了兼容已有逻辑，可以使用CopyF替代】
	CopySF(dest, source interface{}) (err error)
	// CopyF 通用的赋值方法
	CopyF(dest, source interface{}) (err error)
}

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

	xcm XConverters
	acv ActualValueMs
	cp  *sync.Pool
}

const (
	OriginCopyField     = "origin"
	FuncCopyFieldPrefix = "func"
)

// NewCopy 初始化默认
func NewCopy(opts ...Option) *xCopy {
	c := &xCopy{
		convert:   true,
		next:      true,
		recursion: true,
		jsonTag:   true,
		xcm:       xcms,
		acv:       acv,
		cp: &sync.Pool{New: func() interface{} {
			return &convertInfo{}
		}},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Clone 克隆
func (c *xCopy) Clone() *xCopy {
	return &xCopy{
		convert:   c.convert,
		next:      c.next,
		recursion: c.recursion,
		jsonTag:   c.jsonTag,
		xcm:       c.xcm,
		acv:       c.acv,
		cp:        c.cp,
	}
}

// SetConvert 是否强转
func (c *xCopy) SetConvert(convert bool) *xCopy {
	cp := c.Clone()
	WithConvert(convert)(cp)
	return cp
}

// SetNext 出错是否继续赋值下一个字段
func (c *xCopy) SetNext(next bool) *xCopy {
	cp := c.Clone()
	WithNext(next)(cp)
	return cp
}

// SetRecursion 是否递归（依赖强转）
func (c *xCopy) SetRecursion(recursion bool) *xCopy {
	cp := c.Clone()
	WithRecursion(recursion)(cp)
	return cp
}

// SetJSONTag 是否读取JSON TAG
func (c *xCopy) SetJSONTag(jsonTag bool) *xCopy {
	cp := c.Clone()
	WithJsonTag(jsonTag)(cp)
	return cp
}

// 重新定义赋值(带上dv的类型)
func (c *xCopy) copySt(dv, sv reflect.Value) (err error) {
	dt := dv.Type()
	var (
		fst    reflect.StructField
		sf     string
		ofn    bool
		origin bool
		se     error
	)
	data := c.cp.Get().(*convertInfo)
	defer c.cp.Put(data)
	// 赋值所有可能的字段
	num := dt.NumField()
	for i := 0; i < num; i++ {
		fst = dt.Field(i)
		sf, ofn, origin = c.parseTag(fst)

		data.df = fst.Name
		data.osf = origin
		data.dv = dv
		data.sv = sv

		// 解析多级字段并且查看是否符合预期
		sf = strings.Trim(sf, ".")
		if sf != "" && strings.Index(sf, ".") > -1 {
			tmpSfs := strings.Split(sf, ".")
			sfs := make([]string, 0, len(tmpSfs))
			sfi := 0
			for i := range tmpSfs {
				sff := strings.TrimSpace(tmpSfs[i])
				if sff != "" {
					sfs = append(sfs, sff)
					sfi++
				}
			}
			if sfi == 0 {
				sf = ""
			} else if sfi == 1 {
				sf = sfs[0]
			} else {
				sf = sfs[sfi-1]
				aok := c.parseMultiField(sfs[:sfi-1], data)
				if !aok && !c.next {
					err = errors.Errorf("赋值字段[%s]赋值失败: 源字段不存在或为nil", fst.Name)
					return
				}
			}
		}
		// 如果是内嵌字段并且没有指定字段
		if sf == "" && fst.Anonymous {
			dat := fst.Type
			dav := dv.FieldByName(fst.Name)
			isNil := false
			if dat.Kind() == reflect.Ptr {
				dat = dat.Elem()
				if dav.IsNil() {
					if !dav.CanSet() {
						continue
					}
					dav = reflect.New(dat)
					isNil = true
				}
				dav = dav.Elem()
			}
			// 如果是结构体
			if dat.Kind() == reflect.Struct {
				se = c.copySt(dav, data.sv)
				if se == nil && isNil {
					dv.FieldByName(fst.Name).Set(dav.Addr())
				}
			}
		} else {
			data.ofn = ofn
			data.sf = strings.TrimSpace(sf)
			se = c.setSf(data)
		}
		if se == nil || c.next {
			continue
		}
		err = errors.Wrapf(se, "赋值字段[%s]赋值失败", fst.Name)
		return
	}
	return
}

// 解析多级
// 支持数组、map、结构体
func (c *xCopy) parseMultiField(sfs []string, data *convertInfo) (aok bool) {
	defer func() {
		if pe := recover(); pe != nil {
			aok = false
			return
		}
	}()
	for _, sf := range sfs {
		kind := data.sv.Kind()
		for kind == reflect.Ptr {
			if data.sv.IsNil() {
				return false
			}
			data.sv = data.sv.Elem()
			kind = data.sv.Kind()
		}
		ach := c.acv[kind]
		if ach == nil {
			return false
		}
		// 递归重置
		data.sf = strings.TrimSpace(sf)
		err := ach(data)
		if err != nil {
			return false
		}
	}
	// 兜底去除指针
	for data.sv.Kind() == reflect.Ptr {
		if data.sv.IsNil() {
			return false
		}
		data.sv = data.sv.Elem()
	}
	return true
}

// 解析赋值字段
// NOTE
// 优先解析copy，然后解析json，因为按照规则和习惯，大部分的字段最终名与json后的字段名一致
func (c *xCopy) parseTag(fst reflect.StructField) (string, bool, bool) {
	// json -
	// copy origin
	sfi := strings.TrimSpace(fst.Tag.Get("copy"))
	sf := ""
	// split一定不为空
	sfs := strings.SplitN(sfi, ",", 2)
	sf = strings.TrimSpace(sfs[0])
	osf := len(sfs) > 1 && strings.TrimSpace(sfs[1]) == OriginCopyField
	if sf == "" && c.jsonTag {
		sf = strings.TrimSpace(strings.SplitN(strings.TrimSpace(fst.Tag.Get("json")), ",", 2)[0])
	}
	ofn := false
	if len(sf) > 0 {
		sfs = strings.SplitN(sf, ":", 2)
		ofn = strings.TrimSpace(sfs[0]) == FuncCopyFieldPrefix
		if ofn && len(sfs) > 1 {
			// 实际方法名称
			sf = strings.TrimSpace(sfs[1])
		}
	}
	return sf, ofn, osf
}

// 抽离函数，过滤不合法的字段
func (c *xCopy) setSf(data *convertInfo) (err error) {
	// 强制转换可能会出现异常
	defer func() {
		if pe := recover(); pe != nil {
			err = fmt.Errorf("赋值失败: [%#v]", pe)
		}
	}()

	// 检车字段是否有效
	if data.sf == "" {
		data.sf = data.df
	}
	if data.df == "" && data.sf == "" || data.sf == "-" {
		return
	}

	dFieldValue := data.dv.FieldByName(data.df)
	// 判断是否结构体字段存在
	if !dFieldValue.IsValid() || !dFieldValue.CanSet() {
		return
	}
	// 源字段检查
	var sFieldValue reflect.Value
	stKind := data.sv.Kind()
	ach := c.acv[stKind]
	if ach != nil {
		err = ach(data)
		if err != nil {
			return err
		}
		sFieldValue = data.sv
	}
	if !sFieldValue.IsValid() {
		err = errors.New("赋值类型必须是struct、map或者array(slice); 源字段不存在")
		return
	}
	data.dv = dFieldValue
	data.sv = sFieldValue
	c.value(data)
	return
}

// NOTE 复杂逻辑不符合预期: 递归结构体、多级指针
// 抽离函数，真正的赋值操作
func (c *xCopy) value(data *convertInfo) bool {
	dt := data.dv.Type()
	st := data.sv.Type()
	// nil返回
	if (st.Kind() == reflect.Ptr || st.Kind() == reflect.Interface) && data.sv.IsNil() {
		return false
	}

	// 判断类型是否一样或者被赋值是interface
	if dt == st || dt.Kind() == reflect.Interface {
		data.dv.Set(data.sv)
		return true
	}

	// 如果赋值是interface
	if st.Kind() == reflect.Interface {
		data.sv = reflect.ValueOf(data.sv.Interface())
		return c.value(data)
	}

	// 是否强制转换
	if !c.convert {
		return false
	}

	// 找到未初始化的指针
	for dt.Kind() == reflect.Ptr && !data.dv.IsNil() {
		if st.Kind() == reflect.Ptr {
			if data.sv.IsNil() {
				return false
			}
			data.sv = data.sv.Elem()
			st = data.sv.Type()
		}
		data.dv = data.dv.Elem()
		dt = data.dv.Type()
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
		if data.dv != nd.sv {
			data.dv.Set(nd.sv)
		}
		return false
	}
	// 指针过后类型还是不同的话，这里做个判断是为了防止死循环
	dt = nd.dv.Type()
	st = data.sv.Type()
	if st != dt {
		// 如果都为数组
		if (dt.Kind() == reflect.Array || dt.Kind() == reflect.Slice) &&
			(st.Kind() == reflect.Array || st.Kind() == reflect.Slice) {
			return c.recursionSlice(data, nd)
		}
		// 结构体
		if dt.Kind() == reflect.Struct && (st.Kind() == reflect.Map || st.Kind() == reflect.Struct) {
			// note 如果该字段与source是同一个，则会死循环
			return c.recursionStruct(data, nd)
		}
		return c.convertValue(data, nd)
	}
	return true
}

// 强制转换: 调用该方法的时候已经不再是指针
func (c *xCopy) convertValue(data, nd *convertInfo) bool {
	dt := nd.dv.Type()
	sv := data.sv
	nsv := nd.sv
	vh := c.xcm.AC(strings.TrimLeft(dt.PkgPath()+"."+dt.Name(), "."))
	if vh != nil {
		nd.sv = data.sv
		sv = vh.Convert(nd)
	}
	nd.dv.Set(sv.Convert(dt))
	if data.dv != nsv {
		data.dv.Set(nsv)
	}
	return true
}

// 不递归进行指针赋值
func (c *xCopy) notRecursion(data *convertInfo) bool {
	dt := data.dv.Type()
	st := data.sv.Type()
	// 赋值指针(不申请内存拷贝)
	if dt.Kind() == reflect.Ptr {
		// 转换赋值为指针
		if st.Kind() != reflect.Ptr {
			sv := reflect.New(st)
			sv.Elem().Set(data.sv)
			data.sv = sv
		}
		if data.sv.IsNil() {
			return false
		}
		// 赋值
		data.dv.Set(reflect.NewAt(dt.Elem(), unsafe.Pointer(data.sv.Pointer())))
		return true
	}
	// 多级指针
	for st.Kind() == reflect.Ptr {
		if data.sv.IsNil() {
			return false
		}
		data.sv = data.sv.Elem()
		st = data.sv.Type()
	}
	if st == dt {
		data.dv.Set(data.sv)
		return true
	}
	return c.convertValue(data, data)
}

// 递归指针赋值
func (c *xCopy) recursionPointer(data *convertInfo) (*convertInfo, bool) {
	dt := data.dv.Type()
	st := data.sv.Type()
	malloc := false

	nd := c.cp.Get().(*convertInfo)
	// 下一步的赋值全是dv
	nd.dv = data.dv
	nd.sv = data.dv
	nd.df = data.df
	nd.sf = data.sf
	nd.ofn = data.ofn
	nd.osf = data.osf

	for {
		if st == dt {
			nd.dv.Set(data.sv)
			return nd, false
		}
		if dt.Kind() != reflect.Ptr && st.Kind() != reflect.Ptr {
			break
		}
		if dt.Kind() == reflect.Ptr {
			if !malloc {
				// 这里要独立出数据(因为后面可能不需要赋值)
				malloc = true
				nd.dv = reflect.New(dt.Elem())
				nd.sv = nd.dv
			} else {
				nd.dv.Set(reflect.New(dt.Elem()))
			}
			nd.dv = reflect.Indirect(nd.dv)
			dt = nd.dv.Type()
		}
		if st.Kind() == reflect.Ptr {
			if data.sv.IsNil() {
				// 没办法赋值
				return nil, false
			}
			data.sv = reflect.Indirect(data.sv)
			st = data.sv.Type()
		}
	}
	return nd, true
}

// 递归赋值数组
func (c *xCopy) recursionSlice(data, nd *convertInfo) bool {
	sl := data.sv.Len()
	if sl == 0 {
		return false
	}
	tdt := nd.dv.Type()
	if tdt.Kind() == reflect.Array {
		dl := nd.dv.Len()
		if dl < sl {
			sl = dl
		}
	} else {
		data.dv.Set(reflect.MakeSlice(tdt, sl, sl))
	}
	ok := sl == 0 // 等于0算成功
	ndv := nd.dv
	sv := nd.sv
	for i := 0; i < sl; i++ {
		nd.dv = ndv.Index(i)
		nd.sv = data.sv.Index(i)
		if !c.value(nd) {
			continue
		}
		ok = true
	}
	if ok && data.dv != sv {
		data.dv.Set(sv)
	}
	return ok
}

// 递归结构体
func (c *xCopy) recursionStruct(data, nd *convertInfo) bool {
	err := c.copySt(nd.dv, data.sv)
	if err != nil {
		panic(errors.Wrap(err, "递归结构体赋值失败"))
	}
	if data.dv != nd.sv {
		data.dv.Set(nd.sv)
	}
	return true
}

// CopySF 为dest在source中存在的字段自动赋值
// 结构体字段赋值函数
// 通用请调用CopyF方法
func (c *xCopy) CopySF(dest, source interface{}) (err error) {
	if source == nil {
		return errors.New("赋值体不存在")
	}
	defer func() {
		if pe := recover(); pe != nil {
			err = fmt.Errorf("赋值失败: [%#v]", pe)
		}
	}()
	// 校验类型
	dv := reflect.ValueOf(dest)
	dt := dv.Type()
	if dt.Kind() != reflect.Ptr {
		err = errors.New("被赋值的结构体必须是指针类型")
		return
	}

	// 真实数据
	dv = dv.Elem()
	if dv.Kind() != reflect.Struct {
		err = errors.New("被赋值的不是结构体")
		return
	}

	sv := reflect.Indirect(reflect.ValueOf(source))
	// 赋值必须是指针类型
	stKind := sv.Kind()
	if stKind != reflect.Struct && stKind != reflect.Map {
		err = errors.New("赋值类型必须是struct或者map")
		return
	}
	err = c.copySt(dv, sv)
	return
}

// CopyF 单个赋值 通用的赋值方法
func (c *xCopy) CopyF(dest, source interface{}) (err error) {
	if source == nil {
		return errors.New("赋值体不存在")
	}
	// 强制转换可能会出现异常
	defer func() {
		if pe := recover(); pe != nil {
			err = fmt.Errorf("单体赋值失败: [%#v]", pe)
		}
	}()
	// 校验类型
	dv := reflect.ValueOf(dest)
	dt := dv.Type()
	if dt.Kind() != reflect.Ptr {
		return errors.New("被赋值的单体必须是指针类型")
	}
	// 真实数据
	dv = dv.Elem()
	sv := reflect.Indirect(reflect.ValueOf(source))

	// 重置数据
	data := c.cp.Get().(*convertInfo)
	defer c.cp.Put(data)
	data.df = ""
	data.sf = ""
	data.dv = dv
	data.sv = sv

	c.value(data)
	return
}
