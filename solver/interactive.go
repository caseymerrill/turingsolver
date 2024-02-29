package solver

func Interactive() *Solver {
	return &Solver{
		name:             "Interactive",
		codeStrategy:     interactiveCodeStrategy,
		verifierStrategy: interactiveVerifierStrategy,
	}
}

func interactiveCodeStrategy(s *Solver, code []int) int {

}

func interactiveVerifierStrategy(s *Solver, verifierIndex int, code []int) int {

}
