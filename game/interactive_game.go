package game

import (
	"fmt"
	"log"
	"strings"

	"github.com/caseymerrill/turingsolver/verifiers"
)

type InteractiveGame struct {
	cards []*verifiers.VerifierCard
}

func NewInteractiveGame(cards []*verifiers.VerifierCard) Game {
	return &InteractiveGame{
		cards: cards,
	}
}

func (g *InteractiveGame) Stats() map[Player]*PlayerMoves {
	log.Fatal("Stats not supported for interactive game")
	return nil
}

func (g *InteractiveGame) GetVerifierCards() []*verifiers.VerifierCard {
	return g.cards
}

func (g *InteractiveGame) AskQuestion(player Player, code []int, verifierIndex int) bool {
	fmt.Printf("Test code: %v agains verifier %v: %v\n", code, verifierIndex+1, g.cards[verifierIndex])
	return getBoolStdin()
}

func (g *InteractiveGame) MakeGuess(player Player, code []int) bool {
	fmt.Printf("Make a guess: %v", code)
	return getBoolStdin()
}

func (g *InteractiveGame) String() string {
	description := ""
	for _, card := range g.cards {
		verifierDescriptions := make([]string, len(card.Verifiers))
		for i, verifier := range card.Verifiers {
			verifierDescriptions[i] = verifier.Description
		}
		description += strings.Join(verifierDescriptions, " | ")
		description += "\n\n"
	}

	return description
}

func (g *InteractiveGame) Rank() [][]Player {
	log.Fatal("Ranking not supported for interactive game")
	return nil
}

func getBoolStdin() bool {
	for {
		fmt.Print("(y/n): ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal(err)
		}

		switch strings.ToLower(input) {
		case "y", "yes", "true", "1":
			return true
		case "n", "no", "false", "0":
			return false
		default:
			fmt.Println("Unrecognized input")
		}
	}
}
