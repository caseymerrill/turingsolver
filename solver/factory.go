package solver

import "strings"

func FromString(solverName string) *Solver {
	switch strings.ToLower(solverName) {
	case "random":
		return Random()
	case "notimplemented":
		return NotImplementedSolver()
	case "best":
		fallthrough
	case "combine":
		return CombineEliminated()
	case "optimistic":
		return Optimistic()
	case "pessimistic":
		return Pessimistic()
	default:
		panic("Unknown solver name: " + solverName)
	}
}
