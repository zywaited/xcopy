package convert

import "reflect"

type Info struct {
	df  string
	sf  string
	ofn bool
	osf bool
	dv  reflect.Value
	sv  reflect.Value
}

func (ci *Info) GetDf() string {
	return ci.df
}

func (ci *Info) SetDf(df string) {
	ci.df = df
}

func (ci *Info) GetSf() string {
	return ci.sf
}

func (ci *Info) SetSf(sf string) {
	ci.sf = sf
}

func (ci *Info) IsOfn() bool {
	return ci.ofn
}

func (ci *Info) SetOfn(ofn bool) {
	ci.ofn = ofn
}

func (ci *Info) IsOsf() bool {
	return ci.osf
}

func (ci *Info) SetOsf(osf bool) {
	ci.osf = osf
}

func (ci *Info) GetDv() reflect.Value {
	return ci.dv
}

func (ci *Info) SetDv(dv reflect.Value) {
	ci.dv = dv
}

func (ci *Info) GetSv() reflect.Value {
	return ci.sv
}

func (ci *Info) SetSv(sv reflect.Value) {
	ci.sv = sv
}
