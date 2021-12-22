package copy

type Copier interface {
	XCopier

	// --- Copier setters ---

	SetConvert(convert bool) Copier
	SetNext(next bool) Copier
	SetRecursion(recursion bool) Copier
	SetJSONTag(jsonTag bool) Copier
}

type singleXCopy struct {
	c *xCopy
}

// SetConvert 是否强转
func (sc *singleXCopy) SetConvert(convert bool) Copier {
	cp := &singleXCopy{c: sc.c.SetConvert(convert)}
	return cp
}

// SetNext 出错是否继续赋值下一个字段
func (sc *singleXCopy) SetNext(next bool) Copier {
	cp := &singleXCopy{c: sc.c.SetNext(next)}
	return cp
}

// SetRecursion 是否递归（依赖强转）
func (sc *singleXCopy) SetRecursion(recursion bool) Copier {
	cp := &singleXCopy{c: sc.c.SetRecursion(recursion)}
	return cp
}

// SetJSONTag 是否读取JSON TAG
func (sc *singleXCopy) SetJSONTag(jsonTag bool) Copier {
	cp := &singleXCopy{c: sc.c.SetJSONTag(jsonTag)}
	return cp
}

// Copy 实际调用xCopy.CopyF
func (sc *singleXCopy) Copy(dest, source interface{}) error {
	return sc.c.Copy(dest, source)
}

func newSingleXCopy() *singleXCopy {
	return &singleXCopy{c: NewCopy()}
}

var defaultXCopier = newSingleXCopy()

func AcSingleInstance() *singleXCopy {
	return defaultXCopier
}
