package game_generator

import (
	"math/rand"
	"sync/atomic"

	"github.com/caseymerrill/turingsolver/game"
	"github.com/caseymerrill/turingsolver/set"
	"github.com/caseymerrill/turingsolver/solver"
	"github.com/caseymerrill/turingsolver/verifiers"
)

var TotalPotentialSolutions atomic.Int32
var GamesGenrated atomic.Int32
var GamesThrownAway atomic.Int32

func GenerateGame(numberOfVerifierCards int) game.Game {
	solutionFinder := solver.Solver{}
	for {
		cards := make([]*verifiers.VerifierCard, 0, numberOfVerifierCards)
		usedCards := set.Make[int]()
		for len(cards) < numberOfVerifierCards {
			nextCard := rand.Intn(len(verifiers.Cards))
			if usedCards.Contains(nextCard) {
				continue
			}

			usedCards.Add(nextCard)
			cards = append(cards, &verifiers.Cards[nextCard])
		}

		// Interactive game used here because it doesn't require solution/code
		possibleGame := game.NewInteractiveGame(cards)
		solutions := solutionFinder.InitialSolutions(possibleGame)
		if len(solutions) > 0 {
			TotalPotentialSolutions.Add(int32(len(solutions)))
			GamesGenrated.Add(1)
			correctSolution := rand.Intn(len(solutions))
			return game.NewAutoGame(cards, solutions[correctSolution].Verifiers, solutions[correctSolution].Code)
		} else {
			GamesThrownAway.Add(1)
		}
	}
}
