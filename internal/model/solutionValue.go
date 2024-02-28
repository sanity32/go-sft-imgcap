package model

import "github.com/sanity32/b64img"

type SolutionValue struct {
	Solved bool
	Value  bool
}

func FindSolutions(descr Descr, hashes []b64img.Hash) (solutions []bool, hasSolution bool) {
	solutions = make([]bool, len(hashes))
	for n, hash := range hashes {
		sol := Solution{Descr: descr, Hash: hash}.ReadValue()
		if !sol.Solved {
			return
		}
		solutions[n] = sol.Value
	}
	hasSolution = true
	return
}
