package examples

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zywaited/xcopy"
)

func testRecursion(t *testing.T) {
	type destSecond struct {
		Name string
		Age  int
	}
	type destFirst struct {
		User destSecond
	}
	dest := destFirst{}
	source := struct {
		User struct {
			Name string
			Age  int
		}
	}{User: struct {
		Name string
		Age  int
	}{Name: "copy", Age: 22}}
	require.Nil(t, xcopy.Copy(&dest, source))
	require.EqualValues(t, dest.User, source.User)
}
