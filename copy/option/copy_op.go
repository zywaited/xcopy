package option

import "github.com/zywaited/xcopy/copy/convert"

type Config struct {
	Convert   bool
	Next      bool
	Recursion bool
	JsonTag   bool
	Xcm       convert.XConverters
	Acv       convert.ActualValuer
}

type Option func(*Config)

func WithConvert(convert bool) Option {
	return func(c *Config) {
		c.Convert = convert
	}
}

func WithNext(next bool) Option {
	return func(c *Config) {
		c.Next = next
	}
}
func WithRecursion(recursion bool) Option {
	return func(c *Config) {
		c.Recursion = recursion
	}
}

func WithJsonTag(jsonTag bool) Option {
	return func(c *Config) {
		c.JsonTag = jsonTag
	}
}

func WithXCM(xcm convert.XConverters) Option {
	return func(c *Config) {
		c.Xcm = xcm
	}
}

func WithACV(acv convert.ActualValuer) Option {
	return func(c *Config) {
		c.Acv = acv
	}
}
