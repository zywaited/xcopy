package examples

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zywaited/xcopy"
)

func testMultiField(t *testing.T) {
	dest := struct {
		Name string `copy:"db.users.0.name"`
	}{}
	source := struct {
		Db struct {
			Users []map[string]string
		}
	}{Db: struct{ Users []map[string]string }{Users: []map[string]string{{"name": "copy multi name"}}}}
	require.Nil(t, xcopy.Copy(&dest, source))
	require.Equal(t, dest.Name, source.Db.Users[0]["name"])
}
