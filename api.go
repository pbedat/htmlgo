/*
## htmlgo

Type safe and modularize way to generate html on server side.
Download the package with `go get -v github.com/theplant/htmlgo` and import the package with `.` gives you simpler code:

	import (
		. "github.com/theplant/htmlgo"
	)

also checkout full API documentation at: https://godoc.org/github.com/theplant/htmlgo
*/
package htmlgo

import (
	"context"
	"io"
)

type HTMLComponent interface {
	MarshalHTML(ctx context.Context, w io.Writer) error
}

type ComponentFunc func(ctx context.Context, w io.Writer) (err error)

func (f ComponentFunc) MarshalHTML(ctx context.Context, w io.Writer) (err error) {
	return f(ctx, w)
}

type MutableAttrHTMLComponent interface {
	HTMLComponent
	SetAttr(k string, v interface{})
}
