package examples

import (
	"testing"
)

func TestCodecov(t *testing.T) {
	testQuickStart(t)
	testBaseConvert(t)
	testBoolConvert(t)
	testTime(t)
	testTagField(t)
	testStructField(t)
	testMapField(t)
	testSpecialField(t)
	testFieldFunction(t)
	testRecursion(t)
	testPtr(t)
	testMultiField(t)
	testCustom(t)
}
