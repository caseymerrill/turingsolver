package solver

import (
	"bufio"
	"fmt"
	"github.com/caseymerrill/turingsolver/game"
	"os"
	"strings"
)

type InteractiveSolver struct{}

func Interactive() game.Player {
	return InteractiveSolver{}
}

func (i InteractiveSolver) GetPlayerName() string {
	return "Interactive"
}

func (i InteractiveSolver) Solve(g game.Game) (correct bool, solution game.Solution) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		if !reader.Scan() {
			fmt.Println("End of input")
			break
		}

		fmt.Println(g)
		fmt.Println("What do? c=code, g=guess")
		input := reader.Text()
		input = strings.TrimSpace(input)
		if input == "c" {
			i.doCode(g, reader)
		} else if input == "g" {
			correct, solution = i.doGuess(g, reader)
		} else {
			fmt.Println("Invalid input")
		}
	}

	return
}

func (i InteractiveSolver) SetProgressCallback(callback game.ProgressCallback) {
	//TODO implement me
	panic("implement me")
}

func (i InteractiveSolver) Clone() game.Player {
	return InteractiveSolver{}
}

func (i InteractiveSolver) doCode(g game.Game, reader *bufio.Scanner) {
	for {
		fmt.Println("Enter code to test, q to go back: ")
		if !reader.Scan() {
			panic("Failed to read code from stdin")
		}

		input := reader.Text()
		input = strings.TrimSpace(input)
		if input == "q" {
			return
		} else if input == "" {
			fmt.Println("Invalid input must provide code")
			continue
		} else if len(input) != 3 {
			fmt.Println("Invalid input must provide 3 digit code")
			continue
		}

		// Convert input to int slice
		var code []int
		for _, c := range input {
			if c < '1' || c > '5' {
				fmt.Println("Invalid input must provide code with digits 1-5")
				break
			}

			code = append(code, int(c-'0'))
		}

		if len(code) != 3 {
			continue
		}

		for counter := 0; counter < 3; counter++ {
			if !i.doVerifier(code, g, reader) {
				break
			}
		}
	}
}

func (i InteractiveSolver) doVerifier(code []int, g game.Game, reader *bufio.Scanner) bool {
	g.AskQuestion(i, code, 0)
	return true
}

func (i InteractiveSolver) doGuess(g game.Game, reader *bufio.Scanner) (bool, game.Solution) {
	return false, game.Solution{}
}
