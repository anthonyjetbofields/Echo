package echo

import "net/url"

type Context interface {
	QueryParams() url.Values
}

type defaultContext struct {
	queryParams url.Values
}

func (c *defaultContext) QueryParams() url.Values {
	return c.queryParams
}

func NewContext(params url.Values) Context {
	return &defaultContext{queryParams: params}
}
