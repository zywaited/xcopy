package examples

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zywaited/xcopy"
)

func testQuickStart(t *testing.T) {
	// 字段同名同类型赋值
	// dest 待赋值的变量
	// source 数据源
	{
		dest := struct {
			Name string
		}{}
		source := struct {
			Name string
		}{Name: "copy start"}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Name, source.Name)
	}
	// note: source 不能传入 nil interface
	{
		dest := struct {
			Name string
		}{}
		var source interface{}
		require.NotNil(t, xcopy.Copy(&dest, source))
	}
	// note: dest 入参必须是可赋值的指针类型值
	// 1：非指针类型，报错：被赋值的单体必须是指针类型
	{
		dest := struct {
			Name string
		}{}
		source := struct {
			Name string
		}{Name: "copy start"}
		require.NotNil(t, xcopy.Copy(dest, source))
	}
	// 2：指针类型
	// 2.1 nil interface
	{
		var dest interface{}
		source := struct {
			Name string
		}{Name: "copy start"}
		require.NotNil(t, xcopy.Copy(dest, source))
	}
	// 2.2 不可初始化，报错：被赋值的单体无法初始化
	type quickStart struct {
		Name string
	}
	{
		var dest *quickStart
		source := struct {
			Name string
		}{Name: "copy start"}
		require.NotNil(t, xcopy.Copy(dest, source))
	}
	// 2.3 可初始化指针类型，当对指针使用地址引用时可以自动初始化
	{
		var dest *quickStart
		source := struct {
			Name string
		}{Name: "copy start"}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Name, source.Name)
	}
	// 无法转换时会报错
	{
		var dest *quickStart
		source := "copy start"
		require.NotNil(t, xcopy.Copy(&dest, source))
	}
}
