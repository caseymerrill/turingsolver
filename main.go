package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"slices"
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
  turingsolver --gen=<number-of-games> [--n-cards=<number-of-cards>] [--profile] [--solver=<solvers>...]
  turingsolver --print-cards
  
 Options:
  -h --help                    Show this screen.
  --print-cards                Print the available verifier cards.
  --interactive                Run the game in interactive mode.
  --gen=<number-of-games>      Generate <number-of-games> games.
  --n-cards=<number-of-cards>  Generate games with <number-of-cards> verifiers.
  --solver=<solvers>           Use indicated solvers.
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

	var solvers []*solver.Solver
	solverOpt := opts["--solver"]
	if solverOpt == nil {
		solvers = []*solver.Solver{solver.FromString("best")}
	} else {
		solvers = make([]*solver.Solver, len(solverOpt.([]string)))
		for i, solverName := range solverOpt.([]string) {
			solvers[i] = solver.FromString(solverName)
		}
	}

	interactive, _ := opts.Bool("--interactive")
	if interactive {
		interactiveGame := createInteractiveGame()
		interactiveSolver := solvers[0]
		interactiveSolver.SetProgressCallback(func(progress string) {
			fmt.Println(progress)
		})
		_, solution := interactiveSolver.Solve(interactiveGame)
		fmt.Println("Solution:", solution)
	}

	winCount := make(map[string]int)
	winCountLock := sync.Mutex{}
	generateGames, _ := opts.Int("--gen")
	if generateGames > 0 {
		nVerifiers, _ := opts.Int("--n-cards")
		if nVerifiers == 0 {
			nVerifiers = 4
		}

		// games := make([]game.Game, generateGames)
		gameWaitGroup := sync.WaitGroup{}
		for i := 0; i < generateGames; i++ {
			gameWaitGroup.Add(1)
			go func() {
				defer gameWaitGroup.Done()
				gameToSolve := game_generator.GenerateGame(nVerifiers)
				solverWaitGroup := sync.WaitGroup{}
				for _, competingSolver := range solvers {
					solverWaitGroup.Add(1)
					go func(competingSolver *solver.Solver) {
						defer solverWaitGroup.Done()
						var singleSolver solver.Solver = *competingSolver
						correct, solution := singleSolver.Solve(gameToSolve)
						if !correct {
							fmt.Println("Solver", competingSolver.GetPlayerName(), "failed to solve:", solution.Code)
						}
					}(competingSolver)
				}

				solverWaitGroup.Wait()

				winCountLock.Lock()
				defer winCountLock.Unlock()
				rank := gameToSolve.Rank()
				if len(rank) > 0 {
					winners := rank[0]
					// If there is a tie across the board, no one gets a win
					if len(winners) != len(solvers) {
						for _, winner := range rank[0] {
							winCount[winner.GetPlayerName()]++
						}
					}
				}
			}()
		}

		gameWaitGroup.Wait()
		type WinCount struct {
			playerName string
			wins       int
		}
		winners := make([]WinCount, 0, len(winCount))
		for player, wins := range winCount {
			winners = append(winners, WinCount{player, wins})
		}

		slices.SortFunc(winners, func(a, b WinCount) int {
			return b.wins - a.wins
		})

		fmt.Println("Games generated:", game_generator.GamesGenrated.Load())
		fmt.Println("Games thrown away:", game_generator.GamesThrownAway.Load())
		fmt.Println("Average possible solutions:", game_generator.TotalPotentialSolutions.Load()/game_generator.GamesGenrated.Load())

		fmt.Println("Winners:")
		for i, winner := range winners {
			fmt.Printf("\t%v: %v with %v wins\n", i+1, winner.playerName, winner.wins)
		}
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
