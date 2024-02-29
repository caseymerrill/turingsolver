package game

import (
	"fmt"
	"slices"

	"github.com/caseymerrill/turingsolver/verifiers"
)

type Game interface {
	fmt.Stringer
	GetVerifierCards() []*verifiers.VerifierCard
	AskQuestion(player Player, code []int, verifier int) bool
	MakeGuess(player Player, code []int) bool
	// Rank returns the players in order of their performance. Ties go into the same slice
	Rank() [][]Player
}

func PrintWinCount(games []Game) {
	winCount := make(map[string]int)
	for _, solvedGame := range games {
		rank := solvedGame.Rank()
		// If there is a tie across the board, no one gets a win
		if len(rank) > 0 {
			for _, winner := range rank[0] {
				winCount[winner.GetPlayerName()]++
			}
		}
	}

	type WinCount struct {
		playerName string
		wins       int
	}
	winners := make([]WinCount, 0, len(winCount))
	for player, wins := range winCount {
		winners = append(winners, WinCount{player, wins})
	}

	slices.SortFunc(winners, func(a, b WinCount) int {
		return b.wins - a.wins
	})

	fmt.Println("Winners:")
	for i, winner := range winners {
		fmt.Printf("\t%v: %v with %v wins\n", i+1, winner.playerName, winner.wins)
	}
}
