package convert

import "reflect"

func init() {
	xcms.Register("bool", NewBoolXConverter(dc))
}

type BoolXConverter struct {
	next XConverter
}

func NewBoolXConverter(next XConverter) *BoolXConverter {
	return &BoolXConverter{next: next}
}

func (bc *BoolXConverter) Convert(data *Info) reflect.Value {
	sv := data.GetSv()
	if bc.next != nil {
		sv = bc.next.Convert(data)
	}
	return reflect.ValueOf(!sv.IsZero())
}
