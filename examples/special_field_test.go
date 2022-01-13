package examples

import (
	"testing"

	"github.com/zywaited/xcopy"

	"github.com/stretchr/testify/require"
)

func testSpecialField(t *testing.T) {
	dest := struct {
		Id int
	}{}
	source := struct {
		ID uint
	}{ID: 1}
	require.Nil(t, xcopy.Copy(&dest, source))
	require.EqualValues(t, dest.Id, source.ID)
}
