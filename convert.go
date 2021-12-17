package xcopy

import (
	"reflect"
	"strings"
)

// 这个是为了在JSON或者转成MAP更改其字段
// TS是用的驼峰

// Field 字符串转换
type Field func(string) string

func base(data interface{}) (reflect.Type, reflect.Value, bool) {
	dv := reflect.ValueOf(data)
	dt := reflect.TypeOf(data)
	if dt.Kind() != reflect.Ptr {
		// 如果不是指针不能赋值
		pdv := reflect.New(dt)
		pdv.Elem().Set(dv)
		dv = pdv
		dt = pdv.Type()
	}
	for dt.Kind() == reflect.Ptr {
		if dv.IsNil() {
			return dt, dv, false
		}
		dt = dt.Elem()
		dv = dv.Elem()
	}
	return dt, dv, true
}

func ToMapWithField(data interface{}, field Field) map[string]interface{} {
	dt, dv, ok := base(data)
	if !ok {
		return make(map[string]interface{})
	}
	// 支持MAP\STRUCT
	if dt.Kind() == reflect.Struct {
		return structToMap(dv, field)
	}
	if dt.Kind() == reflect.Map {
		return mapToMap(dv, field)
	}
	// ARRAY\SLICE自行循环调用
	return make(map[string]interface{})
}

func ToSliceWithField(data interface{}, field Field) []interface{} {
	dt, dv, ok := base(data)
	if !ok {
		return make([]interface{}, 0)
	}
	if dt.Kind() == reflect.Array || dt.Kind() == reflect.Slice {
		return sliceToMap(dv, field)
	}
	return make([]interface{}, 0)
}

func ToWithField(data interface{}, field Field) interface{} {
	_, dv, ok := base(data)
	if !ok {
		return dv
	}
	rs, _ := selectMapFunc(dv, field)
	return rs
}

// 支持方法
func selectMapFunc(dv reflect.Value, field Field) (interface{}, bool) {
	kind := dv.Kind()
	for kind == reflect.Ptr {
		if dv.IsNil() {
			return nil, false
		}
		dv = dv.Elem()
		kind = dv.Kind()
	}
	// 支持MAP\STRUCT\ARRAY
	if kind == reflect.Struct {
		return structToMap(dv, field), true
	}
	if kind == reflect.Map {
		return mapToMap(dv, field), true
	}
	if kind == reflect.Array || kind == reflect.Slice {
		return sliceToMap(dv, field), true
	}
	return nil, false
}

// 结构体转map
func structToMap(dv reflect.Value, field Field) map[string]interface{} {
	fn := dv.NumField()
	mp := make(map[string]interface{})
	if fn == 0 {
		return mp
	}
	dt := dv.Type()
	for i := 0; i < fn; i++ {
		fs := dt.Field(i)
		fv := dv.Field(i)
		if !fv.CanSet() {
			// 私有变量
			continue
		}
		// 检测TAG
		tag := fs.Tag.Get("json")
		fd := fs.Name
		if tag != "" {
			jfd := strings.TrimSpace(strings.SplitN(tag, ",", 2)[0])
			if jfd != "" {
				if jfd == "-" {
					// 忽略
					continue
				}
				fd = jfd
			}
		}
		if field != nil {
			fd = field(fd)
		}
		if tmp, ok := selectMapFunc(fv, field); ok {
			mp[fd] = tmp
			continue
		}
		mp[fd] = fv.Interface()
	}
	return mp
}

// 该函数只是把内部做转换
func sliceToMap(dv reflect.Value, field Field) []interface{} {
	num := dv.Len()
	is := make([]interface{}, num)
	if num == 0 {
		return is
	}
	for i := 0; i < num; i++ {
		iv := dv.Index(i)
		tmp, ok := selectMapFunc(iv, field)
		if ok {
			is[i] = tmp
			continue
		}
		is[i] = iv.Interface()
	}
	return is
}

func mapToMap(dv reflect.Value, field Field) map[string]interface{} {
	it := dv.MapRange()
	mp := make(map[string]interface{})
	for it.Next() {
		f := it.Key()
		if f.Kind() != reflect.String {
			continue
		}
		fd := f.String()
		if field != nil {
			fd = field(fd)
		}
		v := it.Value()
		if tmp, ok := selectMapFunc(v, field); ok {
			mp[fd] = tmp
			continue
		}
		mp[fd] = v.Interface()
	}
	return mp
}

// To 通用转换
func To(data interface{}) interface{} {
	return ToWithField(data, nil)
}

// ToMap 对结构体字段进行转换
// 依赖JSON TAG
// NOTE 请调用方自行避免使用循环依赖, 与JSON保持一致，私有变量忽略
func ToMap(data interface{}) map[string]interface{} {
	return ToMapWithField(data, nil)
}

// ToSlice 对数组进行每个元素进行转换
func ToSlice(data interface{}) []interface{} {
	return ToSliceWithField(data, nil)
}

// 封装了一些通用的中间函数

// Origin 原始结果
func Origin(next Field) Field {
	return func(s string) string {
		if next != nil {
			s = next(s)
		}
		return s
	}
}

// Snake 转下划线
func Snake(next Field) Field {
	return func(s string) string {
		s = ToSnake(s)
		if next != nil {
			s = next(s)
		}
		return s
	}
}

// Came 转驼峰
func Came(next Field) Field {
	return func(s string) string {
		s = ToCame(s)
		if next != nil {
			s = next(s)
		}
		return s
	}
}

// Upper 转大写
func Upper(next Field) Field {
	return func(s string) string {
		s = strings.ToUpper(s)
		if next != nil {
			s = next(s)
		}
		return s
	}
}

// Lower 转小写
func Lower(next Field) Field {
	return func(s string) string {
		s = strings.ToLower(s)
		if next != nil {
			s = next(s)
		}
		return s
	}
}

// LcFirst 首字母小写
func LcFirst(next Field) Field {
	return func(s string) string {
		s = ToLcFirst(s)
		if next != nil {
			s = next(s)
		}
		return s
	}
}

// UcFirst 首字母大写
func UcFirst(next Field) Field {
	return func(s string) string {
		s = ToUcFirst(s)
		if next != nil {
			s = next(s)
		}
		return s
	}
}
