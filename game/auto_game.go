package game

import (
	"fmt"
	"strings"

	"github.com/caseymerrill/turingsolver/verifiers"
)

type AutoGame struct {
	verifierCards  []*verifiers.VerifierCard
	actualVerfiers []*verifiers.Verifier
	actualCode     []int
}

func NewAutoGame(verifierCards []*verifiers.VerifierCard, actualVerifiers []*verifiers.Verifier, actualCode []int) *AutoGame {
	return &AutoGame{
		verifierCards:  verifierCards,
		actualVerfiers: actualVerifiers,
		actualCode:     actualCode,
	}
}

func (g *AutoGame) String() string {
	description := ""
	for cardIndex, card := range g.verifierCards {
		verifierDescriptions := make([]string, len(card.Verifiers))
		for i, verifier := range card.Verifiers {
			verifierDescriptions[i] = verifier.Description
		}
		description += fmt.Sprintf("Verifier %v: %v\n", cardIndex + 1, strings.Join(verifierDescriptions, " | "))
	}

	return description
}

func (g *AutoGame) GetVerifierCards() []*verifiers.VerifierCard {
	return g.verifierCards
}

func (g *AutoGame) AskQuestion(code []int, verifier int) bool {
	return g.actualVerfiers[verifier].Verify(code...)
}

func (g *AutoGame) MakeGuess(code []int) bool {
	if len(code) != len(g.actualCode) {
		return false
	}

	for i := range code {
		if code[i] != g.actualCode[i] {
			return false
		}
	}

	return true
}
