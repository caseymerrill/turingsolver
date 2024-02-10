package game

import (
	"fmt"

	"github.com/caseymerrill/turingsolver/verifiers"
)

type Game interface {
	fmt.Stringer
	GetVerifierCards() []*verifiers.VerifierCard
	AskQuestion(player Player, code []int, verifier int) bool
	MakeGuess(player Player, code []int) bool
	// Rank returns the players in order of their performance. Ties go into the first slice
	Rank() [][]Player
	Stats() map[Player]*PlayerMoves
}
