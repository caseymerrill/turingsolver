package game

/**
import (
	"strings"

	"github.com/caseymerrill/turingsolver/verifiers"
)

type AutoGame struct {
	VerifierCards []verifiers.VerifierCard
	realVerifiers []verifiers.Verifier
	realCode      []int
}

func NewGame(cards []verifiers.VerifierCard, realVerifiers []verifiers.Verifier, realCode []int) *Game {
	return &AutoGame{
		VerifierCards: cards,
		realVerifiers: realVerifiers,
		realCode:      realCode,
	}
}

func (g *Game) String() string {
	description := ""
	for _, card := range g.VerifierCards {
		verifierDescriptions := make([]string, len(card.Verifiers))
		for i, verifier := range card.Verifiers {
			verifierDescriptions[i] = verifier.Description
		}
		description += strings.Join(verifierDescriptions, " | ")
		description += "\n\n"
	}

	return description
}

func (g *Game) RunVerifier(code []int, verifier int) bool {
	return g.realVerifiers[verifier].Verify(code...)
}

func (g *Game) MakeGuess(code []int) bool {
	if len(code) != len(g.realCode) {
		return false
	}

	for i, n := range code {
		if n != g.realCode[i] {
			return false
		}
	}

	return true
}
*/
