package examples

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zywaited/xcopy"
)

type stringFunction struct {
}

func (sf *stringFunction) String() string {
	return "string-function"
}

func (sf stringFunction) ToString() string {
	return "to-string-function"
}

// *** note: 转成整型时注意是否会溢出，溢出则转换失败 ***
func testBaseConvert(t *testing.T) {
	// integer -> integer
	// *** int8/uint8/int16/uint16/int32/uint32/int64/uint64/int/uint 强制转换同理 ***
	{
		dest := struct {
			Version int
		}{}
		source := struct {
			Version uint16
		}{Version: 1}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.EqualValues(t, dest.Version, source.Version)
	}
	// integer -> string
	// *** int8/uint8/int16/uint16/int32/uint32/int64/uint64/int/uint => string 强制转换同理 ***
	{
		dest := struct {
			Version string
		}{}
		source := struct {
			Version int
		}{Version: 1}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Version, strconv.Itoa(source.Version))
	}
	// string -> string

	// string -> integer
	// *** string => int8/uint8/int16/uint16/int32/uint32/int64/uint64/int/uint 强制转换同理 ***
	{
		dest := struct {
			Version int
		}{}
		source := struct {
			Version string
		}{Version: "1"}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, strconv.Itoa(dest.Version), source.Version)
	}

	// function => string
	{
		dest := ""
		source := stringFunction{}
		require.Nil(t, xcopy.Copy(&dest, &source))
		require.Equal(t, dest, source.String())

		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest, source.ToString())
	}
}

// bool
func testBoolConvert(t *testing.T) {
	// false = integer == 0
	{
		dest := struct {
			Bool bool
		}{}
		source := struct {
			Bool int
		}{}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Empty(t, dest.Bool)

		source.Bool = 1 // != 0
		require.Nil(t, xcopy.Copy(&dest, source))
		require.NotEmpty(t, dest.Bool)
	}
	// false = string == ""
	{
		dest := struct {
			Bool bool
		}{}
		source := struct {
			Bool string
		}{}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Empty(t, dest.Bool)

		source.Bool = "bool" // != ""
		require.Nil(t, xcopy.Copy(&dest, source))
		require.NotEmpty(t, dest.Bool)
	}
	// false = bool false
	// false = (Chan, Func, Interface, Map, Ptr, Slice, UnsafePointer) nil 【引用】
	{
		dest := struct {
			Bool bool
		}{}
		source := struct {
			Bool interface{}
		}{}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Empty(t, dest.Bool)

		source.Bool = []int{} // != nil
		require.Nil(t, xcopy.Copy(&dest, source))
		require.NotEmpty(t, dest.Bool)
	}
	// false = (Array, Struct) all field is false
	{
		dest := struct {
			Bool bool
		}{}
		source := struct {
			Bool [2]int
		}{}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Empty(t, dest.Bool)

		source.Bool = [2]int{1} // not all field is false
		require.Nil(t, xcopy.Copy(&dest, source))
		require.NotEmpty(t, dest.Bool)
	}
	{
		dest := struct {
			Bool bool
		}{}
		source := struct {
			Bool struct{ not bool }
		}{}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Empty(t, dest.Bool)

		source.Bool = struct{ not bool }{not: true} // not all field is false
		require.Nil(t, xcopy.Copy(&dest, source))
		require.NotEmpty(t, dest.Bool)
	}
}

// *** 转成time.Time结构体要求整型是合法的时间戳，字符串是合法的时间格式 2006-01-02 15:04:05 ***
type timeFunction struct {
	now   time.Time
	toNow time.Time
}

func (tf *timeFunction) Time() time.Time {
	return tf.now
}

func (tf timeFunction) ToTime() time.Time {
	return tf.toNow
}

func testTime(t *testing.T) {
	// time.Time => integer
	{
		dest := struct {
			Time int64
		}{}
		source := struct {
			Time time.Time
		}{Time: time.Now()}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Time, source.Time.Unix())
	}
	// time.Time => string
	{
		dest := struct {
			Time string
		}{}
		source := struct {
			Time time.Time
		}{Time: time.Now()}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Time, source.Time.Format("2006-01-02 15:04:05"))
	}

	// integer => time.Time
	{
		dest := struct {
			Time time.Time
		}{}
		source := struct {
			Time int64
		}{Time: time.Now().Unix()}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Time.Unix(), source.Time)
	}
	// string => time.Time
	{
		dest := struct {
			Time time.Time
		}{}
		source := struct {
			Time string
		}{Time: time.Now().Format("2006-01-02 15:04:05")}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Time.Format("2006-01-02 15:04:05"), source.Time)
	}
	// function => time.Time
	{
		var dest time.Time
		source := timeFunction{now: time.Now(), toNow: time.Now().Add(time.Hour)}
		require.Nil(t, xcopy.Copy(&dest, &source))
		require.Equal(t, dest.Unix(), source.Time().Unix())

		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Unix(), source.ToTime().Unix())
	}
}
