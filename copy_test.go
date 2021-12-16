package xcopy

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var c *xCopy

func testStruct(t *testing.T) {
	id := 1
	name := "m"
	{
		type (
			dest struct {
				Id   int
				Name string
			}
			source struct {
				Id   int
				Name *string
			}
		)
		d := dest{}
		s := source{id, &name}
		require.Nil(t, c.CopySF(&d, s))
		require.Equal(t, id, d.Id)
		require.Equal(t, name, d.Name)
	}
	{
		type (
			dest struct {
				Id   int8
				Name *string
			}
			source struct {
				Id   int64
				Name string
			}
		)
		d := dest{}
		s := source{int64(id), name}
		require.Nil(t, c.CopySF(&d, s))
		require.Equal(t, int8(id), d.Id)
		require.NotNil(t, d.Name)
		require.Equal(t, name, *d.Name)
	}
	{
		type (
			dest struct {
				Id int64
			}
			source struct {
				Id *int
			}
		)
		d := dest{}
		s := source{&id}
		require.Nil(t, c.CopySF(&d, s))
		require.Equal(t, int64(id), d.Id)
	}
	{
		type (
			dest struct {
				Id *int8
			}
			source struct {
				Id *int
			}
		)
		d := dest{}
		s := source{&id}
		require.Nil(t, c.CopySF(&d, s))
		require.NotNil(t, d.Id)
		require.Equal(t, int8(id), *d.Id)
	}
	{
		type (
			dest struct {
				Id *int64
			}
			source struct {
				Id int
			}
		)
		d := dest{}
		s := source{id}
		require.Nil(t, c.CopySF(&d, s))
		require.NotNil(t, d.Id)
		require.Equal(t, int64(id), *d.Id)
	}
	{
		type (
			dest struct {
				Id   *int64 `copy:"pid"`
				Name string // `copy:"noname"`
				Age  int    `copy:"real_age"`
			}
			source struct {
				Pid     int
				Name    string
				RealAge int
			}
		)
		d := dest{}
		s := source{Pid: id, Name: "med", RealAge: 18}
		require.Nil(t, c.CopySF(&d, s))
		require.NotNil(t, d.Id)
		require.Equal(t, int64(id), *d.Id)
		require.Equal(t, s.Name, d.Name)
		require.Equal(t, s.RealAge, d.Age)
	}
}

func testMap(t *testing.T) {
	{
		type (
			dest struct {
				Id       int64   `copy:"ID, origin"`
				RealName *string `json:"name"`
				Type     *int8   `copy:", origin "`
				RealAge  int     `copy:"RealAge"`
			}
		)

		st := int32(1)
		var source = map[string]interface{}{
			"ID":       st,
			"name":     "med",
			"Type":     &st,
			"real_age": 18,
		}
		d := &dest{}
		require.Nil(t, c.CopySF(d, source))
		require.EqualValues(t, st, d.Id)
		require.NotNil(t, d.RealName)
		require.NotNil(t, d.Type)
		require.Equal(t, *d.RealName, source["name"])
		require.EqualValues(t, *d.Type, *source["Type"].(*int32))
		require.Equal(t, d.RealAge, source["real_age"])
	}
}

func testRecursion(t *testing.T) {
	// 测试递归指针、map和结构体
	{
		type repeat struct {
			Id int
		}
		type repeats struct {
			Id int64
		}
		type dest struct {
			Id        int8
			DestOne   *dest
			DestTwo   *dest
			DestThree *dest
			Name      **string
			Age       *int
			Alias     [2]int8
			Real      []repeat
		}

		type source struct {
			Id        int64
			DestOne   dest
			DestTwo   map[string]interface{}
			DestThree *source
			Name      *string
			Age       **int
			Alias     [4]int32
			Real      []repeats
		}
		id := 1
		sid := &id
		name := "med"
		s := source{
			Id: int64(id),
			DestOne: dest{
				Id: int8(id) + 1,
				DestOne: &dest{
					Id: int8(id) + 2,
				},
			},
			DestTwo: map[string]interface{}{
				"id": id + 3,
			},
			Name:  &name,
			Age:   &sid,
			Alias: [4]int32{66, 88, 100},
		}
		s.Real = append(s.Real, repeats{1})
		s.Real = append(s.Real, repeats{2})

		d := dest{}
		require.Nil(t, c.CopyF(&d, s))
		require.EqualValues(t, d.Id, s.Id)
		require.NotNil(t, d.DestOne)
		require.Equal(t, *d.DestOne, s.DestOne)
		require.NotNil(t, d.DestTwo)
		require.EqualValues(t, d.DestTwo.Id, s.DestTwo["id"])
		require.NotNil(t, d.Name)
		require.NotNil(t, *d.Name)
		require.Equal(t, **d.Name, *s.Name)
		require.NotNil(t, d.Age)
		require.Equal(t, *d.Age, **s.Age)
		require.EqualValues(t, d.Alias[0], s.Alias[0])
		require.EqualValues(t, d.Alias[1], s.Alias[1])
		require.Equal(t, len(s.Real), len(d.Real))
		require.EqualValues(t, s.Real[0].Id, d.Real[0].Id)
		require.EqualValues(t, s.Real[1].Id, d.Real[1].Id)
	}
}

func testMultiField(t *testing.T) {
	{
		type dest struct {
			private int
			ignore  int    `copy:"-"`
			Id      int    `copy:"ids.id"`
			Name    string `copy:"names.0.name"`
			Age     int    `copy:"ages.f.age"`
			Test    *int   `copy:"t.t.0"`
		}
		type id struct {
			Id int
		}
		type name struct {
			Name string
		}
		type age struct {
			Age int
		}
		type tt struct {
			T []int
		}
		type source struct {
			Private int
			Ignore  int
			Ids     id
			Names   []*name
			Ages    map[string]age
			T       tt
		}
		d := dest{}
		s := source{
			Private: 1,
			Ignore:  1,
			Ids:     id{Id: 1},
			Names:   []*name{{Name: "med"}},
			Ages:    map[string]age{"f": {Age: 5}},
			T:       tt{T: []int{1}},
		}
		c = c.SetNext(false)
		require.Nil(t, c.CopySF(&d, s))
		require.Equal(t, 0, d.private)
		require.Equal(t, 0, d.ignore)
		require.Equal(t, s.Ids.Id, d.Id)
		require.Equal(t, s.Names[0].Name, d.Name)
		require.Equal(t, s.Ages["f"].Age, d.Age)
		require.NotNil(t, d.Test)
		require.Equal(t, s.T.T[0], *d.Test)

		d = dest{}
		require.Nil(t, c.CopyF(&d, s))
		require.Equal(t, 0, d.private)
		require.Equal(t, 0, d.ignore)
		require.Equal(t, s.Ids.Id, d.Id)
		require.Equal(t, s.Names[0].Name, d.Name)
		require.Equal(t, s.Ages["f"].Age, d.Age)
		require.NotNil(t, d.Test)
		require.Equal(t, s.T.T[0], *d.Test)
	}
}

func testAnonymous(t *testing.T) {
	type As struct {
		Id   int
		name string
	}
	type as struct {
		Id   int
		name string
	}
	type d1 struct {
		*As
		Age int
	}
	type source struct {
		as
		Age int
	}
	s := source{
		as:  as{Id: 1},
		Age: 18,
	}
	td := &d1{}
	require.Nil(t, c.CopySF(td, s))
	require.Equal(t, s.Id, td.Id)
	require.Equal(t, s.Age, td.Age)

	type d2 struct {
		as
		Age int
	}
	td2 := &d2{}
	require.Nil(t, c.CopySF(td2, s))
	require.Equal(t, s.Id, td2.Id)
	require.Equal(t, s.Age, td2.Age)
}

func testTimeTo(t *testing.T) {
	type dest struct {
		Now  int
		Next string
	}
	type source struct {
		Now  *time.Time
		Next time.Time
	}
	now := time.Now()
	s := &source{Now: &now, Next: now}
	d := &dest{}
	require.Nil(t, c.CopySF(d, s))
	require.EqualValues(t, d.Now, now.Unix())
	require.EqualValues(t, d.Next, now.Format("2006-01-02 15:04:05"))
}

func testToTime(t *testing.T) {
	type dest struct {
		Now  *time.Time
		Next time.Time
	}
	type source struct {
		Now  int64
		Next string
	}
	now := time.Now()
	d := &dest{}
	s := &source{Now: now.Unix(), Next: now.Format("2006-01-02 15:04:05")}
	require.Nil(t, c.CopySF(d, s))
	require.NotNil(t, d.Now)
	require.EqualValues(t, now.Unix(), (*d.Now).Unix())
	require.EqualValues(t, now.Format("2006-01-02 15:04:05"), d.Next.Format("2006-01-02 15:04:05"))
}

type testCallMethod struct {
}

func (tcm *testCallMethod) String() string {
	return "test-call-method-string"
}

func (tcm *testCallMethod) GetB() int {
	return 200
}

func testMethod(t *testing.T) {
	type dest struct {
		A string
		B int

		C string
		D int
	}

	type source struct {
		A *testCallMethod
		B *testCallMethod

		C int
		D string
	}

	d := &dest{}
	tcm := &testCallMethod{}
	s := &source{
		A: tcm,
		B: tcm,
		C: 100,
		D: "300",
	}
	require.Nil(t, c.CopySF(d, s))
	require.Equal(t, d.A, s.A.String())
	require.Equal(t, d.B, s.B.GetB())
	require.Equal(t, d.C, strconv.Itoa(s.C))
	require.Equal(t, strconv.Itoa(d.D), s.D)
}

func TestCopy(t *testing.T) {
	c = NewCopy()
	testStruct(t)
	testMap(t)
	testRecursion(t)
	testMultiField(t)
	testAnonymous(t)
	testTimeTo(t)
	testToTime(t)
	testMethod(t)
}
