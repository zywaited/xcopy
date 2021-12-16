package xcopy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToCame(t *testing.T) {
	require.Equal(t, "TestId", ToCame("test_id"))
	require.Equal(t, "TestId", ToCame("test_Id"))
	require.Equal(t, "TestId", ToCame("Test_Id"))
	require.Equal(t, "TestId", ToCame("TestId"))
	require.Equal(t, "TestId", ToCame("testId"))
	require.Equal(t, "Testid", ToCame("testid"))
	require.Equal(t, "", ToCame(""))
	require.Equal(t, "", ToCame("_"))
}

func TestToSnake(t *testing.T) {
	require.Equal(t, "test_id", ToSnake("TestId"))
	require.Equal(t, "test_id_", ToSnake("TestId_"))
	require.Equal(t, "_test_id", ToSnake("_testId"))
	require.Equal(t, "testid", ToSnake("testid"))
	require.Equal(t, "", ToSnake(""))
	require.Equal(t, "t", ToSnake("T"))
}
