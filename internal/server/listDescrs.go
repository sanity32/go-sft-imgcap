package server

import (
	"fmt"
	"html/template"
	"net/url"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/sanity32/go-sft-imgcap/internal/model"
)

func (serv *Server) DescrList(ctx *fiber.Ctx) error {
	ll := model.AllDescr()
	if ll == nil {
		return ctx.JSON(map[string]any{"error": "cannot load descriptions"})
	}
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

	var html string
	for _, l := range ll {
		html += fmt.Sprintf("<a href='/descr/distinct/%v'> %v </a><hr />", l, l)
	}
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return ctx.SendString(html)
}

func (serv *Server) DescrDistinct(ctx *fiber.Ctx) error {
	descrStr, _ := url.QueryUnescape(ctx.Params("descr"))
	descr := model.Descr(descrStr)
	if descr == "" {
		return ctx.SendString("ctx not specified")
	}

	fmt.Println("descr:", descr)

	m := descr.Distinct()
	fmt.Printf("m: %#v\n", m)
	var maxValue int
	var highValue int
	items := []ModelDescrDistinctItem{}
	for img, score := range m {
		// b64 := string(b64img.PREFIX_B64_JPG + img)
		// b64 := img.WithJpgPrefix().String()
		b64 := img.String()
		fmt.Println("b64:", b64)
		items = append(items, ModelDescrDistinctItem{
			// B64:       img.WithJpgPrefix().String(),
			MaxValue:  &maxValue,
			HighValue: &highValue,
			Score:     score,
			B64:       b64,
			Hash:      string(img.Hash()),
		})
		if score > maxValue {
			maxValue = score
			highValue = int(float32(maxValue) * 0.75)
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Score > items[j].Score
	})

	funcMap := template.FuncMap{
		"safe": func(s string) template.URL {
			return template.URL(s)
		},
	}

	t, err := template.New("t1").Funcs(funcMap).Parse(`
	<table>
	{{ range .}}
		<tr>
		<td><img src="{{.B64|safe}}"></td>
		<td><meter value={{.Score}} max={{.MaxValue}} high={{.HighValue}} min=0 /></td>
		<td> {{.Hash}} </td>
		<td></td>
		</tr>
	{{ end }}
	</table>
	`)
	if err != nil {
		return ctx.SendString(err.Error())
	}
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return t.Execute(ctx, items)
}
