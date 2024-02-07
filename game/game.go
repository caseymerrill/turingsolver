package game

import (
	"fmt"

	"github.com/caseymerrill/turingsolver/verifiers"
)

type Game interface {
	fmt.Stringer
	GetVerifierCards() []*verifiers.VerifierCard
	AskQuestion(code []int, verifier int) bool
	MakeGuess(code []int) bool
}
