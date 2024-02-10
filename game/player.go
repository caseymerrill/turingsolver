package game

import (
	"slices"

	"github.com/caseymerrill/turingsolver/optional"
	"github.com/caseymerrill/turingsolver/verifiers"
)

const questionsPerCode = 3

type Player interface {
	// GetName returns the name of the player
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

func (p *PlayerMoves) askedQuestion(code []int, card *verifiers.VerifierCard) {
	if p.guessedCorrectly.HasValue() {
		panic("Illegal move player has already guessed")
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
}

func (p *PlayerMoves) madeGuess(code []int, correct bool) {
	if p.guessedCorrectly.HasValue() {
		panic("Illegal move player has already guessed")
	}


	p.codeGuessed = code
	p.guessedCorrectly.Set(correct)
}
