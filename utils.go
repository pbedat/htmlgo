package htmlgo

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
)

type RawHTML string

func (s RawHTML) MarshalHTML(ctx context.Context, w io.Writer) (err error) {
	_, err = fmt.Fprint(w, s)
	return
}

func Text(text string) (r HTMLComponent) {
	return RawHTML(html.EscapeString(text))
}

func Textf(format string, a ...interface{}) (r HTMLComponent) {
	return Text(fmt.Sprintf(format, a...))
}

type HTMLComponents []HTMLComponent

func Components(comps ...HTMLComponent) HTMLComponents {
	return HTMLComponents(comps)
}

func (hcs HTMLComponents) MarshalHTML(ctx context.Context, w io.Writer) (err error) {
	for _, h := range hcs {
		if h == nil {
			continue
		}
		err = h.MarshalHTML(ctx, w)
		if err != nil {
			return
		}
	}
	return
}

func Fprint(w io.Writer, root HTMLComponent, ctx context.Context) (err error) {
	if root == nil {
		return
	}
	err = root.MarshalHTML(ctx, w)

	return
}

func MustString(root HTMLComponent, ctx context.Context) string {
	b := bytes.NewBuffer(nil)
	err := Fprint(b, root, ctx)
	if err != nil {
		panic(err)
	}
	return b.String()
}
