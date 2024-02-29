package solver

import "slices"

// sumBiggestThree returns the sum of the three largest numbers in the slice, will sort the slice int the process
func sumBiggestThree(scores []int) int {
	slices.Sort(scores)
	finalScore := 0
	for _, score := range scores[len(scores)-3:] {
		finalScore += score
	}
	return finalScore
}

// topThreeVerifiersCodeStrategy returns code whose top three verifiers have the highest score
func topThreeVerifiersCodeStategy(s *Solver, code []int) int {
	scores := make([]int, len(s.game.GetVerifierCards()))
	for i := range s.game.GetVerifierCards() {
		scores[i] = s.verifierStrategy(s, i, code)
	}

	return sumBiggestThree(scores)
}
