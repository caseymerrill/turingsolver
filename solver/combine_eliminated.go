package solver

import "slices"

func CombineEliminated() *Solver {
	return &Solver{
		name:            "CombineEliminated",
		codeStrategy:    scoreCodeWithMostEliminatedSolutions,
		verifierStategy: scoreVerifierWithMostEliminations,
	}
}

// scoreCodeWithMostEliminatedSolutions returns the code that could potentially eliminate the most cases (without goign to zero), based on the top three verifiers
func scoreCodeWithMostEliminatedSolutions(s *Solver, code []int) int {
	verifierScores := make([]int, len(s.game.GetVerifierCards()))
	for i := range s.game.GetVerifierCards() {
		verifierScore := 0

		ifValidSolutionCount := len(s.adjustSolutions(code, i, true))
		if ifValidSolutionCount > 0 {
			verifierScore += len(s.solutions) - len(s.adjustSolutions(code, i, true))
		}

		ifInvalidSolutionCount := len(s.adjustSolutions(code, i, false))
		if ifInvalidSolutionCount > 0 {
			verifierScore += len(s.solutions) - len(s.adjustSolutions(code, i, false))
		}

		verifierScores[i] = verifierScore
	}

	// Add up the top three verifier scores since that is how many questions we get per code
	slices.Sort(verifierScores)
	codeScore := 0
	for _, verifierScore := range verifierScores[len(verifierScores)-3:] {
		codeScore += verifierScore
	}

	return codeScore
}

// scoreVerifierWithMostEliminations
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
