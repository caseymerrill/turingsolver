package game

import (
	"fmt"
	"strings"

	"github.com/caseymerrill/turingsolver/verifiers"
)

type Solution struct {
	Code      []int
	Verifiers []*verifiers.Verifier
}

func (sol Solution) String() string {
	code := ""
	for _, c := range sol.Code {
		code += fmt.Sprintf("%v", c)
	}

	verifierDescriptions := make([]string, len(sol.Verifiers))
	for i, verifier := range sol.Verifiers {
		verifierDescriptions[i] = verifier.Description
	}

	return "Code: " + code + ", Verifiers: " + strings.Join(verifierDescriptions, " | ")
}
