package server

import "github.com/gofiber/fiber/v2"

type Server struct {
	App *fiber.App
}

func (serv *Server) simpleOut(ctx *fiber.Ctx, err error) error {
	if err != nil {
		return ctx.JSON(map[string]string{"error": err.Error()})
	}
	return ctx.JSON(map[string]bool{"success": true})

}

func (serv *Server) Listen(addr string) error {
	return serv.App.Listen(addr)
}

func New() *Server {
	serv := Server{}
	app := fiber.New()
	app.Post("/q", serv.postQuery)
	app.Post("/record", serv.postRecord)
	return &Server{App: app}
}
