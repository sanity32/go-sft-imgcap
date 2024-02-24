package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sanity32/go-sft-imgcap/internal/model"
)

type QueryRequest struct {
	model.SessionRecord
	UpdateRecords bool `json:"upd"`
}

type QueryResponse struct {
	Description string    `json:"description"`
	Answers     []bool    `json:"answers"`
	Weights     []float32 `json:"weights"`
}

func (serv *Server) postQuery(ctx *fiber.Ctx) error {
	// var rec model.SessionRecord
	var rec QueryRequest
	if err := ctx.BodyParser(&rec); err != nil {
		return serv.simpleOut(ctx, err)
	}

	if rec.UpdateRecords {
		rec.Write()
		fmt.Println("==== UPDATED ===")
	}

	pool := model.NewNormalPool(rec.Descr)
	if !pool.EnoughSamples() {
		return serv.simpleOut(ctx, model.ErrNotEnoughSamples)
	}

	ww, err := pool.Weights(rec.Hashes())
	if err != nil {
		return serv.simpleOut(ctx, err)
	}

	resp := QueryResponse{
		Description: rec.Descr.String(),
		Weights:     ww,
	}
	resp.Answers = make([]bool, len(ww))
	for n, w := range ww {
		resp.Answers[n] = w >= 1
	}
	return ctx.JSON(resp)
}

func (serv *Server) postRecord(ctx *fiber.Ctx) error {
	var rec model.SessionRecord
	if err := ctx.BodyParser(&rec); err != nil {
		return serv.simpleOut(ctx, err)
	}
	if err := rec.Write(); err != nil {
		return serv.simpleOut(ctx, err)
	}
	return serv.simpleOut(ctx, rec.PopulateHashDir())
}
