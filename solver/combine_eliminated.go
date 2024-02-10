package solver

func CombineEliminated() *Solver {
	return &Solver{
		name:            "CombineEliminated",
		codeStrategy:    topThreeVerifiersCodeStategy,
		verifierStategy: scoreVerifierWithMostEliminations,
	}
}

// scoreVerifierWithMostEliminations scores verifier based off how many solutions are eliminated in all cases.
func scoreVerifierWithMostEliminations(s *Solver, verifierIndex int, code []int) int {
	score := 0
	solutionsIfTrue := len(s.adjustSolutions(code, verifierIndex, true))
	if solutionsIfTrue > 0 {
		score += len(s.solutions) - solutionsIfTrue
	}

	solutionsIfFalse := len(s.adjustSolutions(code, verifierIndex, false))
	if solutionsIfFalse > 0 {
		score += len(s.solutions) - solutionsIfFalse
	}

	return score
}
