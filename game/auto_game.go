package game

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/caseymerrill/turingsolver/verifiers"
)

type AutoGame struct {
	verifierCards   []*verifiers.VerifierCard
	actualVerfiers  []*verifiers.Verifier
	actualCode      []int
	playerStats     map[Player]*PlayerMoves
	playerStatsLock sync.Mutex
}

func NewAutoGame(verifierCards []*verifiers.VerifierCard, actualVerifiers []*verifiers.Verifier, actualCode []int) *AutoGame {
	return &AutoGame{
		verifierCards:  verifierCards,
		actualVerfiers: actualVerifiers,
		actualCode:     actualCode,
		playerStats:    make(map[Player]*PlayerMoves),
	}
}

func (g *AutoGame) String() string {
	description := ""
	for cardIndex, card := range g.verifierCards {
		verifierDescriptions := make([]string, len(card.Verifiers))
		for i, verifier := range card.Verifiers {
			verifierDescriptions[i] = verifier.Description
		}
		description += fmt.Sprintf("Verifier %v: %v\n", cardIndex+1, strings.Join(verifierDescriptions, " | "))
	}

	return description
}

func (g *AutoGame) GetVerifierCards() []*verifiers.VerifierCard {
	return g.verifierCards
}

func (g *AutoGame) AskQuestion(player Player, code []int, verifier int) bool {
	g.playerStatsLock.Lock()
	defer g.playerStatsLock.Unlock()

	playerStats := g.playerStats[player]
	if playerStats == nil {
		playerStats = &PlayerMoves{player: player}
		g.playerStats[player] = playerStats
	}

	if err := playerStats.askedQuestion(code, g.verifierCards[verifier]); err != nil {
		fmt.Println(err)
		return false
	}

	return g.actualVerfiers[verifier].Verify(code...)
}

func (g *AutoGame) MakeGuess(player Player, code []int) (correct bool) {
	g.playerStatsLock.Lock()
	defer g.playerStatsLock.Unlock()

	playerStats := g.playerStats[player]
	if playerStats == nil {
		playerStats = &PlayerMoves{player: player}
		g.playerStats[player] = playerStats
	}

	correct = true
	if len(code) != len(g.actualCode) {
		correct = false
	} else {
		for i := range code {
			if code[i] != g.actualCode[i] {
				correct = false
				break
			}
		}
	}

	if err := playerStats.madeGuess(code, correct); err != nil {
		fmt.Println(err)
		return false
	}

	return
}

func (g *AutoGame) Rank() [][]Player {
	g.playerStatsLock.Lock()
	defer g.playerStatsLock.Unlock()

	stats := make([]*PlayerMoves, 0, len(g.playerStats))
	for _, playerStat := range g.playerStats {
		if playerStat.guessedCorrectly.Value() {
			stats = append(stats, playerStat)
		}
	}
	slices.SortStableFunc(stats, sorter)

	players := make([][]Player, 0, len(stats))
	currentRank := make([]Player, 0, 1)
	if len(stats) == 0 {
		return players
	}

	previousStat := stats[0]
	for _, stat := range stats {
		if stat.codesTested == previousStat.codesTested && len(stat.questionsAsked) == len(previousStat.questionsAsked) {
			currentRank = append(currentRank, stat.player)
		} else {
			players = append(players, currentRank)
			currentRank = []Player{stat.player}
		}

		previousStat = stat
	}

	players = append(players, currentRank)

	return players
}

// sorter assumes all players answered correcty
func sorter(a, b *PlayerMoves) int {
	return (a.codesTested-b.codesTested)*1_000_000 + (len(a.questionsAsked) - len(b.questionsAsked))
}
