package server

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
)

var LastResponse QueryResponse

func (serv Server) RenderLastSuccessful(ctx *fiber.Ctx) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

	var html template.HTML
	q := LastResponse
	html += template.HTML(q.Render())

	return ctx.Send([]byte(html))
}
