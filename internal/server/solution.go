package server

import (
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/sanity32/b64img"
	"github.com/sanity32/go-sft-imgcap/internal/model"
	"github.com/sanity32/go-sft-imgcap/web/ffs"
)

type DescrPostSolutionRequest struct {
	Solved bool `json:"solved"`
	Value  bool `json:"value"`
}

func (serv *Server) DescrSolveList(ctx *fiber.Ctx) error {
	ll := model.AllDescr()
	if ll == nil {
		return ctx.JSON(map[string]any{"error": "cannot load descriptions"})
	}
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

	var html string
	for _, l := range ll {
		html += fmt.Sprintf("<a href='/descr/solve/%v'> %v </a><hr />", l, l)
	}
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return ctx.SendString(html)
}

func (serv *Server) DescrSolve(ctx *fiber.Ctx) error {

	descrStr, _ := url.QueryUnescape(ctx.Params("descr"))
	descr := model.Descr(descrStr)
	if descr == "" {
		return ctx.Status(400).SendString("ctx not specified")
	}

	m := descr.Distinct()
	a := []ffs.Item{}
	for img, _ := range m {
		id := string(img.Hash())
		sol := model.Solution{Descr: descr, Hash: img.Hash()}.ReadValue()
		a = append(a, ffs.Item{
			Id:     id,
			Method: "post",
			Action: "/descr/solve/" + descr.String() + "/" + id + "/",
			ImgSrc: img.WithJpgPrefix().String(),
			Solved: sol.Solved,
			Value:  sol.Value,
		})
	}
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

	return ffs.Render(a, ctx)
}

func (serv *Server) DecrPostSolution(ctx *fiber.Ctx) error {

	descrStr, _ := url.QueryUnescape(ctx.Params("descr"))
	descr := model.Descr(descrStr)
	if descr == "" {
		return ctx.Status(400).SendString("descr not specified")
	}

	hash := ctx.Params("hash")
	if hash == "" {
		return ctx.Status(400).SendString("hash  not specified")
	}

	var req DescrPostSolutionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	sol := model.Solution{Descr: descr, Hash: b64img.Hash(hash)}
	if !req.Solved {
		sol.Delete()
		return ctx.SendStatus(fiber.StatusOK)
	}

	switch req.Value {
	case true:
		sol.Write("1")
	case false:
		sol.Write("0")
	}
	return ctx.JSON(req)
}
