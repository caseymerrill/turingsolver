package game

import (
	"fmt"
	"slices"

	"github.com/caseymerrill/turingsolver/optional"
	"github.com/caseymerrill/turingsolver/verifiers"
)

const questionsPerCode = 3

type Player interface {
	// GetPlayerName returns the name of the player
	GetPlayerName() string
}

type PlayerMoves struct {
	player                 Player
	codesTested            int
	questionsAskedThisCode int
	questionsAsked         []Question

	codeGuessed      []int
	guessedCorrectly optional.Optional[bool]
}

type Question struct {
	Code []int
	Card *verifiers.VerifierCard
}

func (p *PlayerMoves) askedQuestion(code []int, card *verifiers.VerifierCard) error {
	if p.guessedCorrectly.HasValue() {
		return fmt.Errorf("illegal move player has already guessed. Player: ", p.player.GetPlayerName(), " Code: ", code, " Card: ", card)
	}

	var lastCode []int
	if len(p.questionsAsked) > 0 {
		lastCode = p.questionsAsked[len(p.questionsAsked)-1].Code
	}

	if !slices.Equal(lastCode, code) || p.questionsAskedThisCode == questionsPerCode {
		p.codesTested += 1
		p.questionsAskedThisCode = 0
	}

	p.questionsAsked = append(p.questionsAsked, Question{Code: code, Card: card})
	p.questionsAskedThisCode += 1

	return nil
}

func (p *PlayerMoves) madeGuess(code []int, correct bool) error {
	if p.guessedCorrectly.HasValue() {
		return fmt.Errorf("Illegal move player has already guessed. Player: ", p.player.GetPlayerName(), " Code: ", code)
	}

	p.codeGuessed = code
	p.guessedCorrectly.Set(correct)

	return nil
}
