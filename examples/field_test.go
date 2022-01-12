package examples

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zywaited/xcopy"
)

func testTagField(t *testing.T) {
	{
		dest := struct {
			Name string `copy:"Alise"`
		}{}
		source := struct {
			Alise string
		}{Alise: "copy alise name"}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Name, source.Alise)
	}
	{
		dest := struct {
			Name string `json:"Alise"`
		}{}
		source := struct {
			Alise string
		}{Alise: "copy alise name"}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Name, source.Alise)
	}
}

func testStructField(t *testing.T) {
	{
		dest := struct {
			Name string
		}{}
		source := struct {
			name string
			Name string
		}{name: "copy name", Name: "copy uc name"}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Name, source.Name)
	}
	{
		dest := struct {
			Uc_Name string
		}{}
		source := struct {
			UcName string
		}{UcName: "copy uc name"}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Uc_Name, source.UcName)
	}
}

func testMapField(t *testing.T) {
	{
		dest := struct {
			Name string
		}{}
		source := map[string]string{"Name": "copy name"}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Name, source["Name"])
	}
	{
		dest := struct {
			UcName string
		}{}
		source := map[string]string{"uc_name": "copy uc name"}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.UcName, source["uc_name"])
	}
	{
		dest := struct {
			Name string `copy:"name"`
		}{}
		source := map[string]string{"Name": "copy name"}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.Name, source["Name"])
	}
	{
		dest := struct {
			UcName string `copy:"uc_name"`
		}{}
		source := map[string]string{"UcName": "copy uc name"}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.UcName, source["UcName"])
	}
}

type fieldFuncSource struct {
}

func (ffs *fieldFuncSource) UcName() string {
	return "copy name"
}

func (ffs fieldFuncSource) GetUcName() string {
	return "copy uc name"
}

func testFieldFunction(t *testing.T) {
	{
		dest := struct {
			UcName string `json:"uc_name"`
		}{}
		source := &fieldFuncSource{}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.UcName, source.UcName())
	}
	{
		dest := struct {
			UcName string `json:"uc_name"`
		}{}
		source := fieldFuncSource{}
		require.Nil(t, xcopy.Copy(&dest, source))
		require.Equal(t, dest.UcName, source.GetUcName())
	}
}
