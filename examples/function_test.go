package examples

import (
	"testing"

	"github.com/zywaited/xcopy"

	"github.com/stretchr/testify/require"
)

func testFunction(t *testing.T) {
	{
		dest := struct {
			Name string `copy:"func:name"`
		}{}
		source := &user{name: "copy"}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Name, source.Name())
	}
	{
		dest := struct {
			Name string `copy:"func:user.name"`
		}{}
		source := &factory{}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Name, source.User().Name())
	}
}

type user struct {
	name string
}

func (user *user) Name() string {
	return user.name
}

type factory struct {
}

func (factory *factory) User() *user {
	return &user{name: "copy"}
}
