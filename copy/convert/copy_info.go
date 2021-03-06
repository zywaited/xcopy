package convert

import "reflect"

// special field name【id => ID】
var commonInitialisms = map[string]bool{}

func init() {
	// gorm common alise
	for _, field := range []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"} {
		commonInitialisms[field] = true
	}
}

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
