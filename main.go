package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/caseymerrill/turingsolver/game"
	"github.com/caseymerrill/turingsolver/solver"
	"github.com/caseymerrill/turingsolver/verifiers"
	"github.com/docopt/docopt-go"
)

const docString = `TuringSolver

Usage:
  turingsolver [--interactive]
  turingsolver --print-cards
  
 Options:
  -h --help      Show this screen.
  --print-cards  Print the available verifier cards.
  --interactive  Run the game in interactive mode.`

func main() {
	opts, err := docopt.ParseDoc(docString)
	if err != nil {
		log.Fatal(err)
	}

	printCards, _ := opts.Bool("--print-cards")
	if printCards {
		printAllVerifierCards()
		return
	}

	var game game.Game
	interactive, _ := opts.Bool("--interactive")
	if interactive {
		game = createInteractiveGame()
	} else {
		log.Fatal("non-interactive not implemented")
	}

	fmt.Println(game)

	solver := solver.NewSolver(game, func(progress string) {
		fmt.Println(progress)
	})
	solution := solver.Solve()
	fmt.Println("Solution:", solution)
}

func createInteractiveGame() game.Game {
	const prompt = "Add verifiers (blank to stop, - to remove previous): "
	cards := []verifiers.VerifierCard{}
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		if !reader.Scan() {
			break
		}

		input := reader.Text()
		input = strings.TrimSpace(input)

		if input == "" {
			break
		} else if input == "-" {
			if len(cards) > 0 {
				cards = cards[:len(cards)-1]
			}
			continue
		}

		cardNumber, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Could not recognize 'input' as a number")
			continue
		}

		if len(verifiers.Cards) < cardNumber || cardNumber < 1 {
			fmt.Println("No card with that number")
			continue
		}

		fmt.Println("Adding: ", verifiers.Cards[cardNumber-1])
		cards = append(cards, verifiers.Cards[cardNumber-1])
	}

	if reader.Err() != nil {
		log.Fatal(fmt.Errorf("reading input : %w", reader.Err()))
	}

	return game.NewInteractiveGame(cards)
}

func printAllVerifierCards() {
	for i, card := range verifiers.Cards {
		fmt.Printf("%v: %v\n", i, card)
	}
}
