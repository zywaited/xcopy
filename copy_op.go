package xcopy

type Option func(*xCopy)

func WithConvert(convert bool) Option {
	return func(c *xCopy) {
		c.convert = convert
	}
}

func WithNext(next bool) Option {
	return func(c *xCopy) {
		c.next = next
	}
}
func WithRecursion(recursion bool) Option {
	return func(c *xCopy) {
		c.recursion = recursion
	}
}

func WithJsonTag(jsonTag bool) Option {
	return func(c *xCopy) {
		c.jsonTag = jsonTag
	}
}

func WithXCM(xcm XConverters) Option {
	return func(c *xCopy) {
		c.xcm = xcm
	}
}
