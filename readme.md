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
* 自行注册转换和读取类【时间转换、整型和字符串等】

## 赋值字段优先级
### 查询优先级
* 1：copy指定方法
* 2：copy指定字段 > 字段名称 > 字段同名函数 > Get{字段名称}函数

### 转换优先级
* 本身的类型值 > 字段同名函数 > Get{字段名称}函数

### 兜底类型
* bool
* int 整型和字符串自动转换
* string 整型自动转换，兜底可实现String或者ToString方法
* time.Time 整型和字符串自动转换，兜底可实现Time或者ToTime方法

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
        MF        int    `copy:"MultiF.Id"`          // source中的MultiF下的Id字段
        FuncName  string `copy:"func:GetFuncName"`   // source中的MultiF下的GetFuncName方法，依旧支持origin，默认是转成驼峰
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

func (s *source) GetFuncName() string {
    return "copy-func"
}


var anotherSource = map[string]interface{}{"pid": 1, "alias_name": "zy"}
```