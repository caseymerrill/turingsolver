package solver

import (
	"fmt"
	"log"
	"slices"
	"sync"

	"github.com/caseymerrill/turingsolver/game"
	"github.com/caseymerrill/turingsolver/set"
	"github.com/caseymerrill/turingsolver/verifiers"
)

type Solver struct {
	game             game.Game
	progressCallback ProgressCallback
	solutions        []Solution
	codesTested      int
	questionsAsked   int
}

type ProgressCallback func(string)

func NewSolver(gameToSolve game.Game, progressCallback ProgressCallback) *Solver {
	return &Solver{
		game:             gameToSolve,
		progressCallback: progressCallback,
	}
}

func (s *Solver) Solve() Solution {
	s.solutions = s.InitialSolutions()
	s.progressReport()

	for len(s.solutions) > 0 && !s.hasSolution() {
		code := s.bestCodeToAsk()
		s.codesTested += 1

		for i := 0; i < 3; i++ {
			verifier := s.bestVerifiersToAsk(code)
			if verifier == -1 {
				s.progressCallback("Passing, no useful validators for code")
				break
			}

			valid := s.game.AskQuestion(code, verifier)
			s.questionsAsked += 1
			s.solutions = s.adjustSolutions(code, verifier, valid)
			if s.hasSolution() {
				s.finalReport()
				return s.solutions[0]
			}

			s.progressReport()
		}
	}

	if len(s.solutions) == 0 {
		log.Fatal("No solutions found")
	} else {
		resultReport := fmt.Sprintf("Found solution after testing %v codes and asking %v questions:\n", s.codesTested, s.questionsAsked)
		s.progressCallback(resultReport)
	}

	s.finalReport()
	return s.solutions[0]
}

func (s *Solver) InitialSolutions() []Solution {
	solutions := make(chan Solution, 100)
	wg := sync.WaitGroup{}
	for _, vp := range s.getAllVerifierPermutations() {
		wg.Add(1)
		go func(verifierPermutation []*verifiers.Verifier) {
			defer wg.Done()
			var validCode []int
			for _, code := range possibleCodes {
				if verifyCode(code, verifierPermutation) {
					if validCode != nil {
						// One solution per valid verifier permutation
						validCode = nil
						break
					}

					validCode = code
				}
			}

			if validCode != nil && allValidatorsUseful(verifierPermutation) {
				solutions <- Solution{
					Code:      validCode,
					Verifiers: verifierPermutation,
				}
			}
		}(vp)
	}

	go func() {
		wg.Wait()
		close(solutions)
	}()

	var result []Solution
	for solution := range solutions {
		result = append(result, solution)
	}

	return result
}

func (s *Solver) Score() (int, int) {
	return s.codesTested, s.questionsAsked
}

func (s *Solver) progressReport() {
	progressReport := fmt.Sprintf("Found %v solutions:\n", len(s.solutions))
	for _, solution := range s.solutions {
		progressReport += fmt.Sprintf("  %v\n", solution)
	}
	s.progressCallback(progressReport)
}

func (s *Solver) finalReport() {
	resultReport := fmt.Sprintf("Found solution after testing %v codes and asking %v questions:\n", s.codesTested, s.questionsAsked)
	s.progressCallback(resultReport)
}

// hasSolution returns true if all possible solutions use the same code
func (s *Solver) hasSolution() bool {
	if len(s.solutions) == 0 {
		return false
	}

	code := s.solutions[0].Code
	for _, solution := range s.solutions[1:] {
		for i, c := range code {
			if c != solution.Code[i] {
				return false
			}
		}
	}

	return true
}

func (s *Solver) bestCodeToAsk() []int {
	return s.mostEliminatedCases()
}

// mostEliminatedCases returns the code that could potentially eliminate the most cases (without goign to zero), based on the top three verifiers
func (s *Solver) mostEliminatedCases() []int {
	var bestCode []int
	var bestScore int
	currentSolutaionCount := len(s.solutions)
	for _, code := range possibleCodes {
		verifierScores := make([]int, len(s.game.GetVerifierCards()))
		for i := range s.game.GetVerifierCards() {
			score := 0

			ifValidSolutionCount := len(s.adjustSolutions(code, i, true))
			if ifValidSolutionCount > 0 {
				score += currentSolutaionCount - len(s.adjustSolutions(code, i, true))
			}

			ifInvalidSolutionCount := len(s.adjustSolutions(code, i, false))
			if ifInvalidSolutionCount > 0 {
				score += currentSolutaionCount - len(s.adjustSolutions(code, i, false))
			}

			verifierScores[i] = score
		}

		slices.Sort(verifierScores)
		codeScore := 0
		for _, verifierScore := range verifierScores[len(verifierScores)-3:] {
			codeScore += verifierScore
		}

		if codeScore > bestScore {
			bestScore = codeScore
			bestCode = code
		}
	}

	return bestCode
}

func (s *Solver) firstCode() []int {
	return s.solutions[0].Code
}

func (s *Solver) bestVerifiersToAsk(code []int) int {
	bestVerifierIndex := -1
	bestVerifierScore := 0
	for i := range s.game.GetVerifierCards() {
		score := 0
		ifValidSolutionCount := len(s.adjustSolutions(code, i, true))
		if ifValidSolutionCount > 0 {
			score += len(s.solutions) - ifValidSolutionCount
		}

		ifInvalidSolutionCount := len(s.adjustSolutions(code, i, false))
		if ifInvalidSolutionCount > 0 {
			score += len(s.solutions) - ifInvalidSolutionCount
		}

		if score > bestVerifierScore {
			bestVerifierScore = score
			bestVerifierIndex = i
		}
	}

	return bestVerifierIndex
}

func (s *Solver) adjustSolutions(code []int, verifierIndex int, valid bool) []Solution {
	verifiersToKeep := set.Make[*verifiers.Verifier]()
	for _, verifier := range s.game.GetVerifierCards()[verifierIndex].Verifiers {
		if verifier.Verify(code...) == valid {
			verifiersToKeep.Add(verifier)
		}
	}

	newSolutions := make([]Solution, 0, len(s.solutions))
	for _, solution := range s.solutions {
		if verifiersToKeep.Contains(solution.Verifiers[verifierIndex]) {
			newSolutions = append(newSolutions, solution)
		}
	}

	return newSolutions
}

func allValidatorsUseful(verifierPermutation []*verifiers.Verifier) bool {
	for i := range verifierPermutation {
		holdout := make([]*verifiers.Verifier, len(verifierPermutation)-1)
		copy(holdout, verifierPermutation[:i])
		copy(holdout[i:], verifierPermutation[i+1:])
		validCodes := 0
		for _, code := range possibleCodes {
			if verifyCode(code, holdout) {
				validCodes += 1
			}

			if validCodes > 1 {
				// The withheld verifier is useful
				break
			}
		}

		if validCodes <= 1 {
			return false
		}
	}

	return true
}

func verifyCode(code []int, verifierPermutation []*verifiers.Verifier) bool {
	for _, verifier := range verifierPermutation {
		if !verifier.Verify(code...) {
			return false
		}
	}

	return true
}

func join[T fmt.Stringer](left [][]T, right []T) [][]T {
	joined := make([][]T, 0, len(left)*len(right))
	for _, leftOriginal := range left {
		for _, r := range right {
			l := make([]T, len(leftOriginal))
			copy(l, leftOriginal)
			joined = append(joined, append(l, r))
		}
	}

	return joined
}

func (s *Solver) getAllVerifierPermutations() [][]*verifiers.Verifier {
	cards := s.game.GetVerifierCards()
	verifierPermutations := make([][]*verifiers.Verifier, len(cards[0].Verifiers))
	for i, verifier := range cards[0].Verifiers {
		verifierPermutations[i] = []*verifiers.Verifier{verifier}
	}

	for _, card := range cards[1:] {
		verifierPermutations = join(verifierPermutations, card.Verifiers)
	}

	return verifierPermutations
}
