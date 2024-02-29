package solver

import (
	"fmt"
	"slices"

	"github.com/caseymerrill/turingsolver/game"
	"github.com/caseymerrill/turingsolver/set"
	"github.com/caseymerrill/turingsolver/verifiers"
	"gonum.org/v1/gonum/stat/combin"
)

func Combinator() *Solver {
	return &Solver{
		name:             "Combinator",
		codeStrategy:     pessimisticCombinatorCodeStrategy,
		verifierStrategy: pessimisticVerifierStrategy,
	}
}

func Combinator1_1() *Solver {
	return &Solver{
		name:             "Combinator1.1",
		codeStrategy:     combinator1_1CodeStrategy,
		verifierStrategy: combinator1_1VerifierStrategy,
	}
}

func Combinator2() *Solver {
	panic("Not implemented")
	return &Solver{
		name: "Combinator2",
	}
}

func pessimisticCombinatorCodeStrategy(s *Solver, code []int) int {
	neededVerifiers := unsolvedVerifiers(s)
	choose := min(3, len(neededVerifiers))
	combinations := combin.Combinations(len(neededVerifiers), choose)
	bestScore := 0
	for _, combination := range combinations {
		score := 0
		for _, neededVerifierIndex := range combination {
			score += s.verifierStrategy(s, neededVerifiers[neededVerifierIndex], code)
		}

		if score > bestScore {
			bestScore = score
		}
	}

	return bestScore
}

func combinator1_1CodeStrategy(s *Solver, code []int) int {
	combinations := combin.Combinations(len(s.game.GetVerifierCards()), 3)
	bestScore := 0
	for _, combination := range combinations {
		score := 0
		for _, neededVerifierIndex := range combination {
			score += s.verifierStrategy(s, neededVerifierIndex, code)
		}

		if score > bestScore {
			bestScore = score
		}
	}

	return bestScore
}

func scorer(solutions []game.Solution) int {
	solutionCount := len(solutions)
	codeCount := countCodes(solutions)

	return codeCount*1_000_000 + solutionCount
}

func combinator1_1VerifierStrategy(s *Solver, verifierIndex int, code []int) int {
	trueScore := scorer(s.adjustSolutions(code, verifierIndex, true))
	falseScore := scorer(s.adjustSolutions(code, verifierIndex, false))

	worstScore := max(trueScore, falseScore)
	if worstScore == 0 {
		return 0
	}

	return scorer(s.solutions) - worstScore
}

func combinator2(s *Solver, code []int) (int, []int) {
	type possibility struct {
		score            int
		possibleCodes    set.Set[string]
		verifiersIndexes []int
	}

	// currentSolutionCodeCount := countCodes(s.solutions)
	currentCodeCount := countCodes(s.solutions)
	neededVerifiers := unsolvedVerifiers(s)
	choose := min(3, len(neededVerifiers))
	combinations := combin.Permutations(len(neededVerifiers), choose)
	var bestWorstCase possibility
	for _, combination := range combinations {
		resultsForThisCombination := make([]possibility, 0)
		for _, neededVerifierIndex := range combination {
			verifierIndex := neededVerifiers[neededVerifierIndex]
			trueSolutions := s.adjustSolutions(code, verifierIndex, true)
			trueCodes := solutions2CodeSet(trueSolutions)

			falseSolutions := s.adjustSolutions(code, verifierIndex, false)
			falseCodes := solutions2CodeSet(falseSolutions)

			if len(resultsForThisCombination) == 0 {
				if len(trueCodes) > 0 {
					resultsForThisCombination = append(resultsForThisCombination, possibility{possibleCodes: trueCodes, verifiersIndexes: []int{verifierIndex}})
				}

				if len(falseCodes) > 0 {
					resultsForThisCombination = append(resultsForThisCombination, possibility{possibleCodes: falseCodes, verifiersIndexes: []int{verifierIndex}})
				}
			} else {
				newResultPermutations := make([]possibility, 0)
				for _, resultPermutation := range resultsForThisCombination {
					if len(resultPermutation.possibleCodes) == 1 {
						continue
					}

					ifTrue := resultPermutation.possibleCodes.Intersection(trueCodes)
					if len(ifTrue) > 0 {
						newResultPermutations = append(newResultPermutations, possibility{possibleCodes: ifTrue, verifiersIndexes: append(resultPermutation.verifiersIndexes, verifierIndex)})
					}

					ifFalse := resultPermutation.possibleCodes.Intersection(falseCodes)
					if len(ifFalse) > 0 {
						newResultPermutations = append(newResultPermutations, possibility{possibleCodes: ifFalse, verifiersIndexes: append(resultPermutation.verifiersIndexes, verifierIndex)})
					}

					// Keep old possibility even if the next verifier isn't useful
					if len(ifTrue) == 0 && len(ifFalse) == 0 {
						newResultPermutations = append(newResultPermutations, resultPermutation)
					}
				}
				resultsForThisCombination = newResultPermutations
			}
		}

		if len(resultsForThisCombination) == 0 {
			continue
		}

		worstCaseThisCombination := slices.MinFunc(resultsForThisCombination, func(s1, s2 possibility) int {
			return len(s1.possibleCodes) - len(s2.possibleCodes)
		})

		if len(bestWorstCase.possibleCodes) == 0 || len(worstCaseThisCombination.possibleCodes) < len(bestWorstCase.possibleCodes) {
			bestWorstCase = worstCaseThisCombination
		}

	}

	if len(bestWorstCase.possibleCodes) == 0 {
		return -1, nil
	}

	score := currentCodeCount - len(bestWorstCase.possibleCodes)
	//fmt.Printf("%v %v %+v\n", code, score, bestWorstCase)
	return score, bestWorstCase.verifiersIndexes
}

func unsolvedVerifiers(s *Solver) []int {
	codes := set.Set[string]{}
	unsolved := set.Set[int]{}
	firstVerifier := make([]*verifiers.Verifier, len(s.solutions[0].Verifiers))
	for _, solution := range s.solutions {
		// Skip duplicate codes... don't care about solving the verifiers for the same code
		solutionCode := code2str(solution.Code)
		if codes.Contains(solutionCode) {
			continue
		} else {
			codes.Add(solutionCode)
		}

		for verifierIndex, verifier := range solution.Verifiers {
			if firstVerifier[verifierIndex] == nil {
				firstVerifier[verifierIndex] = verifier
			} else if verifier != firstVerifier[verifierIndex] {
				unsolved.Add(verifierIndex)
			}
		}
	}

	asSlice := unsolved.ToSlice()
	slices.Sort(asSlice)
	return asSlice
}

func code2str(code []int) string {
	codeStr := ""
	for _, c := range code {
		codeStr += fmt.Sprintf("%v", c)
	}
	return codeStr
}

func solutions2CodeSet(solutions []game.Solution) set.Set[string] {
	codes := set.Set[string]{}
	for _, solution := range solutions {
		codes.Add(code2str(solution.Code))
	}

	return codes
}

func countCodes(solutions []game.Solution) int {
	return len(solutions2CodeSet(solutions))
}
