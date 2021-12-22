package xcopy

import "github.com/zywaited/xcopy/copy"

// SetConvert 是否强转
func SetConvert(convert bool) copy.Copier {
	return copy.AcSingleInstance().SetConvert(convert)
}

// SetNext 出错是否继续赋值下一个字段
func SetNext(next bool) copy.Copier {
	return copy.AcSingleInstance().SetNext(next)
}

// SetRecursion 是否递归（依赖强转）
func SetRecursion(recursion bool) copy.Copier {
	return copy.AcSingleInstance().SetRecursion(recursion)
}

// SetJSONTag 是否读取JSON TAG
func SetJSONTag(jsonTag bool) copy.Copier {
	return copy.AcSingleInstance().SetJSONTag(jsonTag)
}

// Copy 赋值，实际调用singleXCopy.Copy
func Copy(dest, source interface{}) error {
	return copy.AcSingleInstance().Copy(dest, source)
}
