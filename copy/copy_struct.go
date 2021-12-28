package copy

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/zywaited/xcopy/copy/convert"

	"github.com/pkg/errors"
)

func init() {
	kindCopiers[reflect.Struct] = newStructCopier()
}

func newStructCopier() *structCopier {
	return &structCopier{}
}

type structCopier struct {
}

// 重新定义赋值(带上dv的类型)
func (sc *structCopier) copy(c *xCopy, dv, sv reflect.Value) (err error) {
	dt := dv.Type()
	var (
		fst    reflect.StructField
		sf     string
		ofn    bool
		origin bool
		se     error
	)
	data := c.cp.Get().(*convert.Info)
	defer c.cp.Put(data)
	// 赋值所有可能的字段
	num := dt.NumField()
	for i := 0; i < num; i++ {
		fst = dt.Field(i)
		sf, ofn, origin = sc.parseTag(c, fst)

		data.SetDf(fst.Name)
		data.SetOsf(origin)
		data.SetDv(dv)
		data.SetSv(sv)

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
				aok := sc.parseMultiField(c, sfs[:sfi-1], data)
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
				se = sc.copy(c, dav, data.GetSv())
				if se == nil && isNil {
					dv.FieldByName(fst.Name).Set(dav.Addr())
				}
			}
		} else {
			data.SetOfn(ofn)
			data.SetSf(strings.TrimSpace(sf))
			se = sc.setSf(c, data)
		}
		if se == nil || c.next {
			continue
		}
		err = errors.Wrapf(se, "赋值字段[%s]赋值失败", fst.Name)
		return
	}
	return
}

// 抽离函数，过滤不合法的字段
func (sc *structCopier) setSf(c *xCopy, data *convert.Info) (err error) {
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

// 解析赋值字段
// NOTE
// 优先解析copy，然后解析json，因为按照规则和习惯，大部分的字段最终名与json后的字段名一致
func (sc *structCopier) parseTag(c *xCopy, fst reflect.StructField) (string, bool, bool) {
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

// 解析多级
// 支持数组、map、结构体
func (sc *structCopier) parseMultiField(c *xCopy, sfs []string, data *convert.Info) (aok bool) {
	defer func() {
		if pe := recover(); pe != nil {
			aok = false
			return
		}
	}()
	for _, sf := range sfs {
		kind := data.GetSv().Kind()
		for kind == reflect.Ptr {
			if data.GetSv().IsNil() {
				return false
			}
			data.SetSv(data.GetSv().Elem())
			kind = data.GetSv().Kind()
		}
		ach := c.acv.AC(kind)
		// 递归重置
		data.SetSf(strings.TrimSpace(sf))
		err := ach(data)
		if err != nil {
			return false
		}
	}
	// 兜底去除指针
	for data.GetSv().Kind() == reflect.Ptr {
		if data.GetSv().IsNil() {
			return false
		}
		data.SetSv(data.GetSv().Elem())
	}
	return true
}

// 递归结构体
func (sc *structCopier) recursion(c *xCopy, data, nd *convert.Info) bool {
	err := sc.copy(c, nd.GetDv(), data.GetSv())
	if err != nil {
		panic(errors.Wrap(err, "递归结构体赋值失败"))
	}
	if data.GetDv() != nd.GetSv() {
		data.GetDv().Set(nd.GetSv())
	}
	return true
}

func (sc *structCopier) check(dv, sv reflect.Value) error {
	if reflect.Indirect(dv).Kind() != reflect.Struct {
		return errors.New("被赋值的不是结构体")
	}

	st := reflect.Indirect(sv).Type()
	if st.Kind() != reflect.Map && st.Kind() != reflect.Struct {
		return errors.New("赋值体不是结构体或者哈希")
	}

	return nil
}
