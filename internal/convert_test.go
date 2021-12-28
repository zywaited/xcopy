package internal

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToMap(t *testing.T) {
	type m struct {
		Id   int
		name string
		RA   int `json:"r_a"`
		Am   *m
		M    map[string]int
		Ms   []*m
	}

	m1 := &m{
		Id:   1,
		name: "med",
		RA:   18,
		Am: &m{
			Id:   2,
			name: "med-d2d",
			RA:   19,
			M: map[string]int{
				"Test_Id": 1,
			},
		},
		Ms: []*m{
			{
				Id: 3,
			},
		},
	}
	r := ToMapWithField(m1, Came(LcFirst(nil)))
	s, err := json.Marshal(r)
	require.Nil(t, err)
	fmt.Println(string(s))
}
