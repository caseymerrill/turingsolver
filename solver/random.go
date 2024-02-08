package solver

import "math/rand"

func Random() *Solver {
	return &Solver{
		name:            "Random",
		codeStrategy:    randomCodeStrategy,
		verifierStategy: randomVerifierStrategy,
	}
}

func randomCodeStrategy(s *Solver, code []int) int {
	for i := range s.game.GetVerifierCards() {
		trueSolutionCount := len(s.adjustSolutions(code, i, true))
		falseSolutionCount := len(s.adjustSolutions(code, i, false))
		if (trueSolutionCount > 0 && trueSolutionCount < len(s.solutions)) ||
		(falseSolutionCount > 0 && falseSolutionCount < len(s.solutions)) {
			return rand.Intn(100) + 1
		}
	}

	// Don't choose codes that are totally useless
	return 0
}

func randomVerifierStrategy(s *Solver, verifierIndex int, code []int) int {
	trueSolutionCount := len(s.adjustSolutions(code, verifierIndex, true))
	falseSolutionCount := len(s.adjustSolutions(code, verifierIndex, false))
	if (trueSolutionCount > 0 && trueSolutionCount < len(s.solutions)) ||
	(falseSolutionCount > 0 && falseSolutionCount < len(s.solutions)) {
		return rand.Intn(100) + 1
	}

	// Don't choose verifiers that are totally useless
	return 0
}