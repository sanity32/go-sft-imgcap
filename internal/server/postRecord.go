package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sanity32/b64img"
	"github.com/sanity32/go-sft-imgcap/internal/model"
)

type QueryRequest struct {
	model.SessionRecord
	UpdateRecords bool `json:"upd"`
}

func (serv *Server) postQuery(ctx *fiber.Ctx) error {
	// var rec model.SessionRecord
	var rec QueryRequest
	if err := ctx.BodyParser(&rec); err != nil {
		return serv.simpleOut(ctx, err)
	}

	if rec.UpdateRecords {
		rec.Write()
		rec.PopulateHashDir()
		fmt.Println("==== UPDATED ===")
	}
	hashes := rec.Hashes()

	solvedAnswers, hasSolution := model.FindSolutions(rec.Descr, hashes)
	if hasSolution {
		return ctx.JSON(QueryResponse{
			Description: rec.Descr.String(),
			Hashes:      hashes,
			Answers:     solvedAnswers,
		})
	}

	pool := model.NewNormalPool(rec.Descr)
	if !pool.EnoughSamples() {
		return serv.simpleOut(ctx, model.ErrNotEnoughSamples)
	}

	ww, err := pool.Weights(hashes)
	if err != nil {
		return serv.simpleOut(ctx, err)
	}

	for n, w := range ww {
		if w > .9 && w < 1.15 {
			fmt.Printf("Weight %v is unclear\n", w)
			return ctx.JSON(map[string]any{
				"error": "weight is unclear",
				"n":     n,
				"ww":    ww,
			})
		}
	}

	resp := QueryResponse{
		Description: rec.Descr.String(),
		Weights:     ww,
	}
	resp.Hashes = make([]b64img.Hash, len(hashes))
	resp.Hashes = hashes
	// for n, hash := range hashes {
	// 	resp.Hashes[n] = hash
	// }

	resp.Answers = make([]bool, len(ww))
	for n, w := range ww {
		if solved := solvedAnswers[n]; solved {
			resp.Answers[n] = true
		} else {
			resp.Answers[n] = w >= 1
		}
	}
	// serv.LastSuccessful.Push(resp)
	LastResponse = resp
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
