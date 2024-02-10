package solver

func Pessimistic() *Solver {
	return &Solver{
		name:            "Pessimistic",
		codeStrategy:    topThreeVerifiersCodeStategy,
		verifierStategy: pessimisticVerifierStrategy,
	}
}

func pessimisticVerifierStrategy(s *Solver, verifierIndex int, code []int) int {
	trueSolutionCount := len(s.adjustSolutions(code, verifierIndex, true))
	falseSolutionCount := len(s.adjustSolutions(code, verifierIndex, false))
	worstSolutionsCount := max(trueSolutionCount, falseSolutionCount)
	if worstSolutionsCount == 0 {
		return 0
	}

	return len(s.solutions) - worstSolutionsCount
}
