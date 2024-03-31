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

	"github.com/caseymerrill/turingsolver/server"

	"github.com/caseymerrill/turingsolver/game"
	"github.com/caseymerrill/turingsolver/game_generator"
	"github.com/caseymerrill/turingsolver/solver"
	"github.com/caseymerrill/turingsolver/verifiers"
	"github.com/docopt/docopt-go"
)

const docString = `TuringSolver

Usage:
  turingsolver --interactive [--solver=<solver>]
  turingsolver --server --gen=<number-of-games> [--n-cards=<number-of-cards> --min-solutions=<min-solutions>]
  turingsolver --gen=<number-of-games> [--n-cards=<number-of-cards> --min-solutions=<min-solutions> --profile] [--solver=<solvers>...]
  turingsolver --remote=<url> [--solver=<solvers>...]
  turingsolver --print-cards
  
 Options:
 -h --help                       Show this screen.
 --print-cards                   Print the available verifier cards.
 --interactive                   Run the game in interactive mode.
--gen=<number-of-games>          Generate <number-of-games> games.
--n-cards=<number-of-cards>      Generate games with <number-of-cards> verifiers.
--min-solutions=<min-solutions>  Generate games with at least <min-solutions> solutions.
--solver=<solvers>               Use indicated solvers.
--profile					     Run with CPU profiler.`

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

	numberOfGamesToGenerate, _ := opts.Int("--gen")
	runServer, _ := opts.Bool("--server")
	remoteAdder, _ := opts.String("--remote")

	minSolutions, _ := opts.Int("--min-solutions")
	if minSolutions == 0 {
		minSolutions = 2
	}

	nVerifiers, _ := opts.Int("--n-cards")
	if nVerifiers == 0 {
		nVerifiers = 4
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
	} else if runServer {
		fmt.Println("Generating Games...")
		games := generateGames(numberOfGamesToGenerate, nVerifiers, minSolutions)
		fmt.Println("Starting Server...")
		gameServer := server.NewGameServer(games)
		gameServer.Listen()
	} else if remoteAdder != "" {
		wg := sync.WaitGroup{}
		for _, solverToUse := range solvers {
			remoteGames, err := game.JoinGames(remoteAdder, solverToUse.GetPlayerName())
			if err != nil {
				log.Fatal("Joining games : ", err)
			}

			wg.Add(1)
			go func(solverToUse *solver.Solver) {
				defer wg.Done()
				for _, remoteGame := range remoteGames {
					correct, _ := solverToUse.Solve(remoteGame)
					if !correct {
						fmt.Println("Solver", solverToUse.GetPlayerName(), "failed to solve")
					}
				}
			}(solverToUse)
		}

		wg.Wait()
	} else if numberOfGamesToGenerate > 0 {
		evaluateSolvers(numberOfGamesToGenerate, nVerifiers, minSolutions, solvers)
	}
}

func evaluateSolvers(numberOfGamesToGenerate int, nVerifiers int, minSolutions int, solvers []*solver.Solver) {
	fmt.Println("Generating Games...")
	games := generateGames(numberOfGamesToGenerate, nVerifiers, minSolutions)
	fmt.Println("Solving...")
	gameWaitGroup := sync.WaitGroup{}
	for _, gameToSolve := range games {
		for _, competingSolver := range solvers {
			gameWaitGroup.Add(1)
			go func(competingSolver *solver.Solver, gameToSolve game.Game) {
				defer gameWaitGroup.Done()
				singleSolver := *competingSolver
				correct, solution := singleSolver.Solve(gameToSolve)
				if !correct {
					fmt.Println("Solver", competingSolver.GetPlayerName(), "failed to solve:", solution.Code)
				}
			}(competingSolver, gameToSolve)
		}
	}

	gameWaitGroup.Wait()
	game.PrintWinCount(games)
}

func generateGames(numberOfGamesToGenerate, nVerifiers, minSolutions int) []game.Game {
	games := make(chan game.Game, numberOfGamesToGenerate/10+1)
	wg := sync.WaitGroup{}

	for i := 0; i < numberOfGamesToGenerate; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			games <- game_generator.GenerateGame(nVerifiers, minSolutions)
		}()

	}

	go func() {
		wg.Wait()
		close(games)
	}()

	result := make([]game.Game, 0, numberOfGamesToGenerate)
	for g := range games {
		result = append(result, g)
	}

	return result
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

		if cardNumber > 1002 { // 1002 is the lowest Extreme code (low number is always first)
			lowNumber := cardNumber / 1000
			highNumber := cardNumber % 1000

			if verifyCardNumber(lowNumber) && verifyCardNumber(highNumber) {
				fmt.Println("Adding XTREAM card:", lowNumber, ":", highNumber)
				xtreamCard := verifiers.Cards[lowNumber-1].Combine(verifiers.Cards[highNumber-1])
				cards = append(cards, &xtreamCard)
			} else {
				fmt.Println("No XTREAM card with that number:", cardNumber, ":", lowNumber, ":", highNumber)
				continue
			}
		} else if verifyCardNumber(cardNumber) {
			fmt.Println("Adding: ", verifiers.Cards[cardNumber-1])
			cards = append(cards, &verifiers.Cards[cardNumber-1])
		} else {
			fmt.Println("No card with that number")
			continue
		}
	}

	if reader.Err() != nil {
		log.Fatal(fmt.Errorf("reading input : %w", reader.Err()))
	}

	return game.NewInteractiveGame(cards)
}

func verifyCardNumber(cardNumber int) bool {
	return cardNumber <= len(verifiers.Cards) && cardNumber > 0
}

func printAllVerifierCards() {
	for i, card := range verifiers.Cards {
		fmt.Printf("%v: %v\n", i, card)
	}
}
