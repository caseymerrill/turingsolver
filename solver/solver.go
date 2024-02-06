package solver

import (
	"fmt"
	"log"
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
	s.solutions = s.initialSolutions()

	progressReport := fmt.Sprintf("Found %v initial solutions:\n", len(s.solutions))
	for _, solution := range s.solutions {
		progressReport += fmt.Sprintf("  %v\n", solution)
	}
	s.progressCallback(progressReport)

	for len(s.solutions) > 0 && !s.hasSolution() {
		code := s.bestCodeToAsk()
		s.codesTested += 1

		for i := 0; i < 3; i++ {
			verifier := s.bestVerifiersToAsk()
			if verifier == -1 {
				break
			}

			valid := s.game.AskQuestion(code, verifier)
			s.questionsAsked += 1
			s.solutions = s.adjustSolutions(code, verifier, valid)
			if s.hasSolution() {
				return s.solutions[0]
			}
		}
	}

	if len(s.solutions) == 0 {
		log.Fatal("No solutions found")
	} else {
		progressReport = fmt.Sprintf("Found solution after testing %v codes and asking %v questions:\n", s.codesTested, s.questionsAsked)
		s.progressCallback(progressReport)
	}

	return s.solutions[0]
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

func (s *Solver) verifierCasesToSolve() []set.Set[*verifiers.Verifier] {
	cards := s.game.GetVerifierCards()
	result := make([]set.Set[*verifiers.Verifier], len(cards))
	for i := range cards {
		result[i] = set.Make[*verifiers.Verifier]()
		for _, solution := range s.solutions {
			result[i].Add(solution.Verifiers[i])
		}
	}

	return result
}

func (s *Solver) bestCodeToAsk() []int {
	return s.solutions[0].Code
}

func (s *Solver) bestVerifiersToAsk() int {
	toSolve := s.verifierCasesToSolve()
	for i := range toSolve {
		if len(toSolve[i]) > 1 {
			return i
		}
	}

	return -1
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
		if !verifiersToKeep.Contains(solution.Verifiers[verifierIndex]) {
			newSolutions = append(newSolutions, solution)
		}
	}

	return newSolutions
}

func (s *Solver) initialSolutions() []Solution {
	solutions := make([]Solution, 0)
	for _, verifierPermutation := range s.getAllVerifierPermutations() {
		var validCode []int
		codes, cleanupCodes := possibleCodes()
		defer cleanupCodes()
		for code := range codes {
			if verifyCode(code, verifierPermutation) {
				if validCode != nil {
					// One solution per valid verifier permutation
					validCode = nil
					break
				}

				validCode = code
			}
		}
		cleanupCodes()

		if validCode != nil && allValidatorsUseful(verifierPermutation) {
			solutions = append(solutions, Solution{
				Code:      validCode,
				Verifiers: verifierPermutation,
			})
		}
	}

	return solutions
}

func allValidatorsUseful(verifierPermutation []*verifiers.Verifier) bool {
	for i := range verifierPermutation {
		holdout := make([]*verifiers.Verifier, len(verifierPermutation)-1)
		copy(holdout, verifierPermutation[:i])
		copy(holdout[i:], verifierPermutation[i+1:])
		codes, cleanup := possibleCodes()
		defer cleanup()
		validCodes := 0
		for code := range codes {
			if verifyCode(code, holdout) {
				validCodes += 1
			}

			if validCodes > 1 {
				// The withheld verifier is useful
				cleanup()
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

type empty struct{}

// possibleCodes returns a channel that will send all possible codes, and cleanup function to end generation.
func possibleCodes() (chan []int, func()) {
	codes := make(chan []int)
	cancel := make(chan empty)
	go func() {
		defer close(codes)
		for i := 1; i < 6; i++ {
			for j := 1; j < 6; j++ {
				for k := 1; k < 6; k++ {
					select {
					case codes <- []int{i, j, k}:
					case <-cancel:
						return
					}
				}
			}
		}
	}()

	return codes, sync.OnceFunc(func() { close(cancel) })
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
