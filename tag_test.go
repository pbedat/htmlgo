package htmlgo_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"text/template"

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

func BenchmarkHtmlgo(b *testing.B) {
	f, err := os.OpenFile("/dev/null", os.O_WRONLY, 0600)
	if err != nil {
		b.Fatal(err)
	}

	w := bufio.NewWriter(f)

	for i := 0; i < b.N; i++ {

		var children HTMLComponents

		for i := 0; i < 200; i++ {
			children = append(children, Li(nested(i), Text(fmt.Sprint(i))))
		}

		Ul(children).MarshalHTML(context.TODO(), w)
		//Fprint(f, Ul(children), context.TODO())
	}
}

func nested(j int) HTMLComponent {

	root := Div(Text("test"))

	var children HTMLComponents

	for i := 0; i < 200; i++ {
		children = append(children, Div(Text("test"), Text(fmt.Sprint(i, j))))
	}

	return root.Children(children)
}

func BenchmarkTemplate(b *testing.B) {
	tpl, _ := template.New("test").Parse(`<ul>{{range $i, $_ := .Items}}<li>
<div>
test
	{{ range $j, $_ := $.Items}}<div>test {{$i}} {{$j}}</div>{{end}}
</div>
</li>{{end}}</ul>`)

	f, err := os.OpenFile("/dev/null", os.O_WRONLY, 0600)
	if err != nil {
		b.Fatal(err)
	}

	w := bufio.NewWriter(f)

	items := make([]byte, 200)

	for i := 0; i < b.N; i++ {
		if err := tpl.Execute(w, struct {
			Items []byte
		}{
			Items: items,
		}); err != nil {
			b.Fatal(err)
		}
	}
}
