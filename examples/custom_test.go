package examples

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zywaited/xcopy"
	copy2 "github.com/zywaited/xcopy/copy"
	"github.com/zywaited/xcopy/copy/convert"
	"github.com/zywaited/xcopy/copy/option"
)

type floatConvert struct {
}

func (fc *floatConvert) Convert(data *convert.Info) reflect.Value {
	sv := data.GetSv()
	if !sv.IsValid() {
		return sv
	}
	if sv.Type().Kind() != reflect.String {
		return sv
	}
	fv, err := strconv.ParseFloat(sv.String(), 64)
	if err != nil {
		return sv
	}
	return reflect.ValueOf(fv)
}

func testCustom(t *testing.T) {
	// note: 全局生效
	xcm := convert.AcDefaultXConverter().(convert.XConvertersSetter)
	xcm.Register("float64", &floatConvert{})
	dest := float64(0)
	source := "1"
	require.Nil(t, xcopy.Copy(&dest, source))
	require.Equal(t, strconv.FormatFloat(dest, 'f', 0, 64), source)

	// 当前生效
	cxcm := convert.AcDefaultXConverter().(convert.XConvertersCloner).Clone()
	cxcm.(convert.XConvertersSetter).Register("float64", &floatConvert{})
	cp := copy2.NewCopy(option.WithXCM(cxcm))
	dest = float64(0)
	require.Nil(t, cp.Copy(&dest, source))
	require.Equal(t, strconv.FormatFloat(dest, 'f', 0, 64), source)
}
