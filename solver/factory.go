package solver

import "strings"

func FromString(solverName string) *Solver {
	switch strings.ToLower(solverName) {
	case "random":
		return Random()
	case "notimplemented":
		return NotImplementedSolver()
	case "combine":
		return CombineEliminated()
	case "optimistic":
		return Optimistic()
	case "pessimistic":
		return Pessimistic()
	case "pc":
		return Combinator()
	case "best":
		fallthrough
	case "pc1.1":
		return Combinator1_1()
	case "pc2":
		return Combinator2()
	default:
		panic("Unknown solver name: " + solverName)
	}
}
