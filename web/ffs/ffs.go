package ffs

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
)

//go:embed ffs.tmpl
var ffsTmpl string

//go:embed ffs.js
var ffsJs string

//go:embed ffs.css
var ffsCss string

type Item struct {
	Id     string
	Action string
	Method string
	ImgSrc string
	Solved bool
	Value  bool
}

func Render(ffs []Item, wr io.Writer) error {
	if len(ffs) < 1 {
		msg := ` [ warning: list is empty ] `
		wr.Write([]byte(msg))
		return nil
	}
	if _, err := wr.Write([]byte(fmt.Sprintf("\n<style>\n%v\n</style>\n", ffsCss))); err != nil {
		return err
	}

	fncs := template.FuncMap{"safe": func(s string) template.URL { return template.URL(s) }}
	z := template.Must(template.New("ffs").Funcs(fncs).Parse(ffsTmpl))
	for _, item := range ffs {
		z.Execute(wr, item)
	}

	if _, err := wr.Write([]byte(fmt.Sprintf(`\n<script>
	%v;
	%v("%v", "%v")
	</script>\n`,
		ffsJs, "makeFfs", ".ffs", "value"))); err != nil {
		return err
	}
	return nil
}
