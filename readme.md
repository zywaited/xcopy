# xcopy
赋值工具，能够有效减少各个协议字段之间的互相

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
### 快速入门
* 字段同名，类型一致
```go
package main

import "github.com/zywaited/xcopy"

func main() {
    // 字段同名同类型赋值
    // dest 待赋值的变量
    // source 数据源
    dest := struct {
        Name string // 只能设置可导出字段
    }{}
    source := struct {
        Name string
    }{Name: "copy start"}
    // 使用
    // dest 本身为可初始化地址或者取地址才能被赋值
    _ = xcopy.Copy(&dest, source)
}
```
#### 错误用法
##### source为空
```go
    dest := struct {
        Name string
    }{}
    var source interface{}
    // err: 赋值体不存在
    xcopy.Copy(&dest, source)
```
##### dest为空
```go
    var dest interface{}
    source := struct {
        Name string
    }{Name: "copy start"}
    // err: 被赋值的单体必须有效
    xcopy.Copy(dest, source)
```
##### dest 非指针
```go
    dest := struct {
        Name string
    }{}
    source := struct {
        Name string
    }{Name: "copy start"}
    // err: 被赋值的单体必须是指针类型
    xcopy.Copy(dest, source)
```
##### dest 不可初始化
```go
    type quickStart {
        Name string
    }
    var dest *quickStart
    source := struct {
        Name string
    }{Name: "copy start"}
    // err: 被赋值的单体无法初始化
    xcopy.Copy(dest, source)
	
    // 正确使用：取地址会自动初始化
    // dest.Name = source.Name
    xcopy.Copy(&dest, source)
```

### 类型转换
* convert为true时生效
#### 基础
##### integer2Integer
* 转成整型时注意是否会溢出，溢出则转换失败
```go
    // *** int8/uint8/int16/uint16/int32/uint32/int64/uint64/int/uint 强制转换同理 ***
    dest := struct {
        Version int
    }{}
    source := struct {
        Version uint16
    }{Version: 1}
    xcopy.Copy(&dest, source)
    // dest.Version == source.Version == 1
```
##### integer2String
```go
    // *** int8/uint8/int16/uint16/int32/uint32/int64/uint64/int/uint => string 强制转换同理 ***
    dest := struct {
        Version string
    }{}
    source := struct {
        Version int
    }{Version: 1}
    xcopy.Copy(&dest, source)
    // dest.Version == strconv.Itoa(source.Version) == "1"
```
##### string2Integer
* 转成整型时注意是否会溢出，溢出则转换失败
```go
    // *** string => int8/uint8/int16/uint16/int32/uint32/int64/uint64/int/uint 强制转换同理 ***
    dest := struct {
        Version int
    }{}
    source := struct {
        Version string
    }{Version: "1"}
    xcopy.Copy(&dest, source)
    // strconv.Itoa(dest.Version) == source.Version == "1"
```

#### 特殊转换
##### 布尔
###### integer
* integer != 0
```go
    dest := struct {
        Bool bool
    }{}
    source := struct {
        Bool int
    }{}
    xcopy.Copy(&dest, source)
    // dest.Bool == false
    
    source.Bool = 1 // != 0
    xcopy.Copy(&dest, source)
    // dest.Bool == true
```
###### string
* string != ""
```go
    dest := struct {
        Bool bool
    }{}
    source := struct {
        Bool string
    }{}
    xcopy.Copy(&dest, source)
    // dest.Bool == false
    
    source.Bool = "bool" // != ""
    xcopy.Copy(&dest, source)
    // dest.Bool == true
```
###### Chan, Func, Interface, Map, Ptr, Slice, UnsafePointer
* (Chan, Func, Interface, Map, Ptr, Slice, UnsafePointer) != nil
```go
    dest := struct {
        Bool bool
    }{}
    source := struct {
        Bool interface{}
    }{}
    xcopy.Copy(&dest, source)
    // dest.Bool == false
    
    source.Bool = []int{} // != nil
    xcopy.Copy(&dest, source)
    // dest.Bool == true
```
###### Array, Struct
* (Array, Struct) one field is true
```go
    {
        dest := struct {
            Bool bool
        }{}
        source := struct {
            Bool [2]int
        }{}
        xcopy.Copy(&dest, source)
        // dest.Bool == false
    
        source.Bool = [2]int{1} // not all field is false
        xcopy.Copy(&dest, source)
        // dest.Bool == true
    }
    {
        dest := struct {
            Bool bool
        }{}
        source := struct {
            Bool struct{ not bool }
        }{}
        xcopy.Copy(&dest, source)
        // dest.Bool == false
    
        source.Bool = struct{ not bool }{not: true} // not all field is false
        xcopy.Copy(&dest, source)
        // dest.Bool == true
    }
}
```

##### 时间【time.Time】
###### time2Integer [秒]
* 转成整型时注意是否会溢出，溢出则转换失败
```go
    dest := struct {
        Time int64
    }{}
    source := struct {
        Time time.Time
    }{Time: time.Now()}
    xcopy.Copy(&dest, source)
    // dest.Time == source.Time.Unix()
```
###### time2String
```go
    dest := struct {
        Time string
    }{}
    source := struct {
        Time time.Time
    }{Time: time.Now()}
    xcopy.Copy(&dest, source)
    // dest.Time == source.Time.Format("2006-01-02 15:04:05")
```
###### integer2Time
* 数字要为时间戳：秒
```go
    dest := struct {
        Time time.Time
    }{}
    source := struct {
        Time int64
    }{Time: time.Now().Unix()}
    xcopy.Copy(&dest, source)
    // dest.Time.Unix() == source.Time
```
###### string2Time
* 字符串格式符合: 2006-01-02 15:04:05
```go
    dest := struct {
        Time time.Time
    }{}
    source := struct {
        Time string
    }{Time: time.Now().Format("2006-01-02 15:04:05")}
    xcopy.Copy(&dest, source)
    // dest.Time.Format("2006-01-02 15:04:05") == source.Time
```
##### 函数转换 [优先级最低]
###### string函数
* String【返回值必须是string】
* ToString【返回值必须是string】
```go
    // defined
    type stringFunction struct {
    }
    func (sf *stringFunction) String() string {
        return "string-function"
    }
    // used
    dest := ""
    source := stringFunction{}
    xcopy.Copy(&dest, &source)
    // dest == source.String()
```
##### 时间【time.Time】
* Time【返回值必须是time.Time】
* ToTime【返回值必须是time.Time】
```go
    // defined
    type timeFunction struct {
        now   time.Time
        toNow time.Time
    }
    func (tf *timeFunction) Time() time.Time {
        return tf.now
    }
    // used
    var dest time.Time
    source := timeFunction{now: time.Now(), toNow: time.Now().Add(time.Hour)}
    xcopy.Copy(&dest, &source)
    // dest.Unix() == source.Time().Unix()
```

### 寻找source赋值字段
* 优先级依次降低
#### copy tag
```go
    dest := struct {
        Name string `copy:"Alise"`
    }{}
    source := struct {
        Alise string
    }{Alise: "copy alise name"}
    xcopy.Copy(&dest, source)
    // dest.Name == source.Alise
```
#### json tag
* 该功能需要jsonTag设置为true
```go
    dest := struct {
        Name string `json:"Alise"`
    }{}
    source := struct {
        Alise string
    }{Alise: "copy alise name"}
    xcopy.Copy(&dest, source)
    // dest.Name == source.Alise
```
#### struct字段
##### 字段首字母大写
```go
    dest := struct {
        Name string
    }{}
    source := struct {
        name string
        Name string
    }{name: "copy name", Name: "copy uc name"}
    xcopy.Copy(&dest, source)
    // dest.Name == source.Name
```
##### 可导出驼峰字段
```go
    dest := struct {
        Uc_Name string
    }{}
    source := struct {
        UcName string
    }{UcName: "copy uc name"}
    xcopy.Copy(&dest, source)
    // dest.Uc_Name == source.UcName
```
#### map字段
##### 同名字段
```go
    dest := struct {
        Name string
    }{}
    source := map[string]string{"Name": "copy name"}
    xcopy.Copy(&dest, source)
    // dest.Name == source["Name"]
```
##### 蛇形字段
```go
    dest := struct {
        UcName string
    }{}
    source := map[string]string{"uc_name": "copy uc name"}
    xcopy.Copy(&dest, source)
    // dest.UcName == source["uc_name"]
```
##### 字段首字母大写
```go
    dest := struct {
        Name string `copy:"name"`
    }{}
    source := map[string]string{"Name": "copy name"}
    xcopy.Copy(&dest, source)
    // dest.Name == source["Name"]
```
##### 可导出驼峰字段
```go
    dest := struct {
        UcName string `copy:"uc_name"`
    }{}
    source := map[string]string{"UcName": "copy uc name"}
    xcopy.Copy(&dest, source)
    // dest.UcName == source["UcName"]
```
#### 字段函数[可导出驼峰函数]
```go
    // defined
    type fieldFuncSource struct {
    }
    func (ffs *fieldFuncSource) UcName() string {
        return "copy name"
    }
    // used
    dest := struct {
        UcName string `json:"uc_name"`
    }{}
    source := &fieldFuncSource{}
    xcopy.Copy(&dest, source)
    // dest.UcName == source.UcName()
```
#### Get{字段驼峰}函数
```go
    // defined
    type fieldFuncSource struct {
    }
    func (ffs *fieldFuncSource) GetUcName() string {
        return "copy name"
    }
    // used
    dest := struct {
        UcName string `json:"uc_name"`
    }{}
    source := &fieldFuncSource{}
    xcopy.Copy(&dest, source)
    // dest.UcName == source.GetUcName()
```
### 函数使用易错点
```go
    // 特殊注意函数的作用类型
    // defined
    type fieldFuncSource struct {
    }
    func (ffs *fieldFuncSource) UcName() string {
        return "copy name"
    }
    func (ffs fieldFuncSource) GetUcName() string {
        return "copy uc name"
    }
    // used
    dest := struct {
        UcName string `json:"uc_name"`
    }{}
    source := fieldFuncSource{}
    xcopy.Copy(&dest, source)
    // dest.UcName == source.GetUcName()
    // 虽然说UcName的优先级高于GetUcName, 但是source用的实际数据，不是地址，所以不能访问UcName函数，所以只能取GetUcName
```