package solver

import (
	"math"
)

func Optimistic() *Solver {
	return &Solver{
		name:            "Optimistic",
		codeStrategy:    topThreeVerifiersCodeStategy,
		verifierStategy: optimisticVerifierStrategy,
	}
}

func optimisticVerifierStrategy(s *Solver, verifierIndex int, code []int) int {
	trueSolutionCount := len(s.adjustSolutions(code, verifierIndex, true))
	falseSolutionCount := len(s.adjustSolutions(code, verifierIndex, false))
	bestSolutionsCount := int(math.Min(float64(trueSolutionCount), float64(falseSolutionCount)))
	if bestSolutionsCount == 0 {
		return 0
	}

	return len(s.solutions) - bestSolutionsCount
}
