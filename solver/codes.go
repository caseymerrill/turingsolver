package solver

import (
	"math"
)

const maxDigit = 5

var possibleCodes [][]int

// possibleCodes returns a channel that will send all possible codes, and cleanup function to end generation.
func init() {
	possibleCodes = make([][]int, 0, int(math.Pow(maxDigit, 3)))
	for i := 1; i <= maxDigit; i++ {
		for j := 1; j <= maxDigit; j++ {
			for k := 1; k <= maxDigit; k++ {
				possibleCodes = append(possibleCodes, []int{i, j, k})
			}
		}
	}
}
