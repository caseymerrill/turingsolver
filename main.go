package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"

	"github.com/caseymerrill/turingsolver/game"
	"github.com/caseymerrill/turingsolver/game_generator"
	"github.com/caseymerrill/turingsolver/solver"
	"github.com/caseymerrill/turingsolver/verifiers"
	"github.com/docopt/docopt-go"
)

const docString = `TuringSolver

Usage:
  turingsolver --interactive [--solver=<solver>]
  turingsolver --gen=<number-of-games> [--n-cards=<number-of-cards>] [--profile] [--solver=<solver>]
  turingsolver --print-cards
  
 Options:
  -h --help                    Show this screen.
  --print-cards                Print the available verifier cards.
  --interactive                Run the game in interactive mode.
  --gen=<number-of-games>      Generate <number-of-games> games.
  --n-cards=<number-of-cards>  Generate games with <number-of-cards> verifiers.
  --solver=<solver>            Use indicated solver
  --profile					   Run with CPU profiler.`

func main() {
	opts, err := docopt.ParseDoc(docString)
	if err != nil {
		log.Fatal(err)
	}

	profile, _ := opts.Bool("--profile")
	if profile {
		profileFilename := os.Args[0] + ".prof"
		fmt.Println("Profiling Enabled. Writing to", profileFilename)
		profileOutput, err := os.Create(profileFilename)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(profileOutput)
		defer pprof.StopCPUProfile()
	}

	printCards, _ := opts.Bool("--print-cards")
	if printCards {
		printAllVerifierCards()
		return
	}

	var solverToUse *solver.Solver
	solverOpt, _ := opts.String("--solver")
	if solverOpt == "" {
		solverToUse = solver.FromString("best")
	} else {
		solverToUse = solver.FromString(solverOpt)
	}

	interactive, _ := opts.Bool("--interactive")
	if interactive {
		interactiveGame := createInteractiveGame()
		solverToUse.SetProgressCallback(func(progress string) {
			fmt.Println(progress)
		})
		solution := solverToUse.Solve(interactiveGame)
		fmt.Println("Solution:", solution)
	}

	generateGames, _ := opts.Int("--gen")
	if generateGames > 0 {
		nVerifiers, _ := opts.Int("--n-cards")
		if nVerifiers == 0 {
			nVerifiers = 4
		}

		// games := make([]game.Game, generateGames)
		wg := sync.WaitGroup{}
		for i := 0; i < generateGames; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				gameToSolve := game_generator.GenerateGame(nVerifiers)
				var singleSolver solver.Solver = *solverToUse
				singleSolver.Solve(gameToSolve)
				// codesTested, questionsAsked := solver.Score()
				// fmt.Printf("Solution: %v\nCodes tested: %v\nQuestions asked: %v\n\n", solution, codesTested, questionsAsked)
			}()
		}

		wg.Wait()
	}
}

func createInteractiveGame() game.Game {
	const prompt = "Add verifiers (blank to stop, - to remove previous): "
	cards := []*verifiers.VerifierCard{}
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
		cards = append(cards, &verifiers.Cards[cardNumber-1])
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
