package model

import (
	"os"
	"path"

	"github.com/sanity32/b64img"
)

type Solution struct {
	Descr Descr
	Hash  b64img.Hash
}

func (sol Solution) Dir() string {
	return path.Join(ASSETS_DIR, "solutions", sol.Descr.String())
}

func (sol Solution) Filepath() string {
	return path.Join(sol.Dir(), string(sol.Hash))
}

func (sol Solution) Read() string {
	if bb, err := os.ReadFile(sol.Filepath()); err != nil {
		return ""
	} else {
		return string(bb)
	}
}

func (sol Solution) Write(v string) error {
	os.MkdirAll(sol.Dir(), 0644)
	return os.WriteFile(sol.Filepath(), []byte(v), 0644)
}

func (sol Solution) Delete() error {
	return os.Remove(sol.Filepath())
}

func (sol Solution) ReadValue() SolutionValue {
	s := sol.Read()
	switch s {
	case "0":
		return SolutionValue{Solved: true, Value: false}
	case "1":
		return SolutionValue{Solved: true, Value: true}
	}
	return SolutionValue{Solved: false}
}

func (sol Solution) WriteValue(v SolutionValue) error {
	if !v.Solved {
		return sol.Delete()
	}
	opts := map[bool]string{
		false: "0",
		true:  "1",
	}
	return sol.Write(opts[v.Value])
}
