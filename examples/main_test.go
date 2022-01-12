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
	testFieldFunction(t)
}
