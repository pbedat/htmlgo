package htmlgo_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	. "github.com/theplant/htmlgo"
	"github.com/theplant/testingutils"
)

var htmltagCases = []struct {
	name     string
	tag      *HTMLTagBuilder
	expected string
}{
	{
		name: "case 1",
		tag: Div(
			Div().Text("Hello"),
		),
		expected: `
<div>
<div>Hello</div>
</div>
`,
	},
	{
		name: "case 2",
		tag: Div(
			Div().Text("Hello").
				Attr("class", "menu",
					"id", "the-menu",
					"style").
				Attr("id", "menu-id"),
		),
		expected: `
<div>
<div class='menu' id='menu-id'>Hello</div>
</div>
`,
	},
	{
		name: "escape 1",
		tag: Div(
			Div().Text("Hello").
				Attr("class", "menu",
					"id", "the><&\"'-menu",
					"style"),
		),
		expected: `
<div>
<div class='menu' id='the><&"&#39;-menu'>Hello</div>
</div>
`,
	},
}

func TestHtmlTag(t *testing.T) {
	for _, c := range htmltagCases {
		buf := bytes.NewBuffer([]byte{})
		err := c.tag.MarshalHTML(context.TODO(), buf)
		if err != nil {
			panic(err)
		}
		diff := testingutils.PrettyJsonDiff(c.expected, buf.String())
		if len(diff) > 0 {
			t.Error(c.name, diff)
		}
	}
}

func BenchmarkList(b *testing.B) {
	f, err := os.OpenFile("/dev/null", os.O_WRONLY, 0600)

	w := bufio.NewWriter(f)

	for i := 0; i < b.N; i++ {

		var children HTMLComponents

		for i := 0; i < 200; i++ {
			children = append(children, Li(nested(i), Text(fmt.Sprint(i))))
		}

		if err != nil {
			b.Fatal(err)
		}

		Ul(children).MarshalHTML(context.TODO(), w)
		//Fprint(f, Ul(children), context.TODO())
	}
}

func nested(j int) HTMLComponent {

	root := Div(Text("test"))

	current := root

	for i := 0; i < 50; i++ {
		next := Div(Text("test"), Text(fmt.Sprint(i, j)))

		current.Children(next)
		current = next
	}

	return root
}
