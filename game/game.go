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
	Rank() []Player
}

type Player interface {
	GetPlayerName() string
	Solve(game Game) Solution
}
