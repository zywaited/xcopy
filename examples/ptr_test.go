package examples

import (
	"testing"

	"github.com/zywaited/xcopy"

	"github.com/stretchr/testify/require"
)

func testPtr(t *testing.T) {
	{
		type ptr struct {
		}
		var dest **ptr
		var sourcePtr *ptr
		source := &sourcePtr // **ptr
		require.Nil(t, xcopy.Copy(&dest, &source))
		require.NotNil(t, dest)
	}
	{
		type ptr struct {
		}
		var dest **ptr
		var source *ptr
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Nil(t, dest)
	}
	{
		type ptr struct {
		}
		var dest **ptr
		var sourcePtr *ptr
		source := &sourcePtr // **ptr
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Nil(t, dest)
	}
	{
		type ptr struct {
		}
		var dest **ptr
		sourcePtr := ptr{}
		source := &sourcePtr // **ptr
		require.Nil(t, xcopy.Copy(&dest, source))
		require.NotNil(t, dest)
	}
}
