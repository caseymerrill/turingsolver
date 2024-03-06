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
	name string

	// codeStrategy used for picking codes to guess
	codeStrategy CodeStrategy

	// verifierStrategy used for picking verifiers to query
	verifierStrategy VerifierStrategy

	// progressCallback will be called with a string describing the progress so far, may be left nil
	progressCallback game.ProgressCallback

	// verifiersTestedThisCode is the number of verifiers tested for the current code
	verifiersTestedThisCode int

	game      game.Game
	solutions []game.Solution
}

type CodeStrategy func(*Solver, []int) int
type VerifierStrategy func(*Solver, int, []int) int
type CombinedStrategy func(*Solver, []int) (int, []int)

// NotImplementedSolver can't solve any games, but can be used to get the inital solutions for game generation
func NotImplementedSolver() *Solver {
	return &Solver{
		name: "Not Implemented",
		codeStrategy: func(*Solver, []int) int {
			panic("Not implemented")
		},
		verifierStrategy: func(*Solver, int, []int) int {
			panic("Not implemented")
		},
	}
}

func (s *Solver) GetPlayerName() string {
	return s.name
}

func (s *Solver) Clone() game.Player {
	s2 := *s
	return &s2
}

func (s *Solver) SetProgressCallback(callback game.ProgressCallback) {
	s.progressCallback = callback
}

// Reset clears game, and solution state, but keep configuration like progress callback
func (s *Solver) reset() {
	s.game = nil
	s.solutions = nil
}

func (s *Solver) Solve(gameToSolve game.Game) (bool, game.Solution) {
	s.solutions = s.InitialSolutions(gameToSolve)
	s.progressReport()

	var codesTested [][]int
	for len(s.solutions) > 0 && !s.hasSolution() {
		code := s.selectCode()

		for i := 0; i < 3; i++ {
			var verifier int
			verifier = s.selectVerifier(code)

			if verifier == -1 {
				if s.progressCallback != nil {
					s.progressCallback("No useful verifiers for code")
				}
				break
			}

			valid := s.game.AskQuestion(s, code, verifier)
			s.verifiersTestedThisCode += 1
			s.solutions = s.adjustSolutions(code, verifier, valid)
			if s.hasSolution() {
				break
			}

			s.progressReport()
		}

		codesTested = append(codesTested, code)
		if len(codesTested) > 100 {
			fmt.Println(s.GetPlayerName(), " has tested 100 codes, giving up")
			fmt.Printf("Codes tested:\n%+v\n", codesTested)
			fmt.Println("Solutions:")
			for _, solution := range s.solutions {
				fmt.Printf("\t%+v\n", solution)
			}
			return false, game.Solution{}
		}
	}

	if len(s.solutions) == 0 {
		log.Fatal("No solutions found")
	}

	return s.game.MakeGuess(s, s.solutions[0].Code), s.solutions[0]
}

func (s *Solver) InitialSolutions(gameToSolve game.Game) []game.Solution {
	s.reset()
	s.game = gameToSolve
	solutions := make(chan game.Solution, 100)
	wg := sync.WaitGroup{}
	for vp := range s.getAllVerifierPermutations() {
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
				solutions <- game.Solution{
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

	var result []game.Solution
	for solution := range solutions {
		result = append(result, solution)
	}

	return result
}

func (s *Solver) progressReport() {
	if s.progressCallback == nil {
		return
	}

	progressReport := fmt.Sprintf("Found %v solutions:\n", len(s.solutions))
	for _, solution := range s.solutions {
		progressReport += fmt.Sprintf("  %v\n", solution)
	}

	s.progressCallback(progressReport)
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

func (s *Solver) selectCode() []int {
	bestScore := -1
	var bestCode []int
	for _, code := range possibleCodes {
		var score int
		score = s.codeStrategy(s, code)

		if score > bestScore {
			bestScore = score
			bestCode = code
		}
	}

	return bestCode
}

func (s *Solver) selectVerifier(code []int) int {
	bestVerifierIndex := -1
	bestVerifierScore := 0
	for i := range s.game.GetVerifierCards() {
		score := s.verifierStrategy(s, i, code)
		if score > bestVerifierScore {
			bestVerifierScore = score
			bestVerifierIndex = i
		}
	}

	return bestVerifierIndex
}

func (s *Solver) adjustSolutions(code []int, verifierIndex int, valid bool) []game.Solution {
	verifiersToKeep := set.Make[*verifiers.Verifier]()
	for _, verifier := range s.game.GetVerifierCards()[verifierIndex].Verifiers {
		if verifier.Verify(code...) == valid {
			verifiersToKeep.Add(verifier)
		}
	}

	newSolutions := make([]game.Solution, 0, len(s.solutions))
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

func (s *Solver) getAllVerifierPermutations() chan []*verifiers.Verifier {
	result := make(chan []*verifiers.Verifier, 100)

	go func() {
		defer close(result)
		s.verifierPermutationHelper([]int{}, result)
	}()

	return result
}

func (s *Solver) verifierPermutationHelper(indexesSoFar []int, result chan []*verifiers.Verifier) {
	if len(indexesSoFar) == len(s.game.GetVerifierCards()) {
		verifierPermutation := make([]*verifiers.Verifier, len(indexesSoFar))
		for i, index := range indexesSoFar {
			verifierPermutation[i] = s.game.GetVerifierCards()[i].Verifiers[index]
		}

		result <- verifierPermutation
		return
	}

	for i := 0; i < len(s.game.GetVerifierCards()[len(indexesSoFar)].Verifiers); i++ {
		s.verifierPermutationHelper(append(indexesSoFar, i), result)
	}
}
