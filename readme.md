# xcopy
赋值工具

## 支持场景
* 整型赋值: int->(int,uint...., string, time)
* 字符串赋值: string->(int,uint..., time)
* 结构体、数组、map互相递归赋值

## 特点
* 递归赋值（根据类型不同自动选择）
* 多级字段指定
* 多指针赋值
* 自行注册转换和读取类

## 实例
其他用法查阅copy_test，实现了各种场景和用法

```go
type (
	dest struct {
		Id        *int64 `copy:"pid"`      // source中Pid
		Name      string `copy:"noname"`   // source不存在noname, 使用Name
		Age       int    `copy:"real_age"` // source中的RealAge
		Ignore    bool   `copy:"-"`        // 忽略该字段
		Status    bool   // source中的Status
		AliasName string `copy:"alias_name, origin"` // origin代表不对copy中的值做转换
        MF        int    `copy:"MultiF.Id"` // source中的MultiF下的Id字段
	}
	source struct {
		Pid     int
		Name    string
		RealAge int
		Ignore  bool
        MultiF  struct {
            Id int
        }
	}
)

var anotherSource = map[string]interface{}{"pid": 1, "alias_name": "med"}
```