package server

import (
	"fmt"
	"html/template"

	"github.com/sanity32/b64img"
	"github.com/sanity32/go-sft-imgcap/internal/model"
)

type QueryResponse struct {
	Description string        `json:"description"`
	Answers     []bool        `json:"answers"`
	Weights     []float32     `json:"weights"`
	Hashes      []b64img.Hash `json:"hashes"`
}

func (q QueryResponse) Render() template.HTML {
	fmt.Printf("Rendering %#v", q)
	var html string
	html += `<div class="query-response-render">`
	html += fmt.Sprintf(`<b class="descr"> %v </b>`, q.Description)
	html += `<table> <tr>`
	for n := range q.Answers {
		html += "<td>"
		img, _ := model.MainHashDir.Read(q.Hashes[n])
		color := "red"
		if q.Answers[n] {
			color = "green"
		}
		html += fmt.Sprintf(`<img src="%v" style="border-color: %v;"/><hr /><font color="%v">%v</font>`, img.String(), color, color, q.Weights[n])
		html += "</td>"
	}
	html += `</tr> </td>`
	html += `</div>`
	return template.HTML(html)
}
