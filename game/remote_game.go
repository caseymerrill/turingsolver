package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/caseymerrill/turingsolver/types"
	"github.com/caseymerrill/turingsolver/verifiers"
)

type RemoteGame struct {
	addr          string
	client        *http.Client
	gameIndex     int
	verifierCards []*verifiers.VerifierCard
}

func JoinGames(addr string, playerName string) ([]Game, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("initializing cookiejar : %w", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	if err := join(addr, client, playerName); err != nil {
		return nil, fmt.Errorf("joining game : %w", err)
	}

	games, err := getGames(addr, client)
	if err != nil {
		return nil, fmt.Errorf("getting games : %w", err)
	}

	remoteGames := make([]Game, len(games))
	for i, game := range games {
		cards := make([]*verifiers.VerifierCard, len(game))
		for j, cardNumber := range game {
			if cardNumber < 1 || cardNumber > len(verifiers.Cards) {
				return nil, fmt.Errorf("invalid card number: %v", cardNumber)
			}

			cards[j] = &verifiers.Cards[cardNumber-1]
		}

		remoteGames[i] = &RemoteGame{
			addr:          addr,
			client:        client,
			gameIndex:     i,
			verifierCards: cards,
		}
	}

	return remoteGames, nil
}

func join(addr string, client *http.Client, playerName string) error {
	request := types.JoinRequest{PlayerName: playerName}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	response, err := client.Post(addr+"/join", "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		return fmt.Errorf("joining game : %w", err)
	} else if err := checkResponse(response); err != nil {
		return err
	}

	return nil
}

func getGames(addr string, client *http.Client) ([][]int, error) {
	response, err := client.Get(addr + "/player/games")
	if err != nil {
		return nil, fmt.Errorf("getting games : %w", err)
	} else if err := checkResponse(response); err != nil {
		return nil, err
	}

	responseBody := types.GetGamesResponse{}
	responseDecoder := json.NewDecoder(response.Body)
	if err := responseDecoder.Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("decoding response : %w", err)
	}

	return responseBody.Games, nil
}

func checkResponse(response *http.Response) error {
	if response.StatusCode != 200 {
		return fmt.Errorf("response status code: %v", response.StatusCode)
	}

	return nil
}

func (g *RemoteGame) String() string {
	return "Remote game: " + g.addr
}

func (g *RemoteGame) GetVerifierCards() []*verifiers.VerifierCard {
	return g.verifierCards
}

func (g *RemoteGame) AskQuestion(player Player, code []int, verifier int) bool {
	request := types.AskQuestionRequest{
		GameIndex:     g.gameIndex,
		VerifierIndex: verifier,
		Code:          code,
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling request while testing verifier:", err)
		return false
	}

	response, err := g.client.Post(g.addr+"/player/test-verifier", "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		fmt.Println("error testing verifier:", err)
		return false
	} else if err := checkResponse(response); err != nil {
		fmt.Println("Response error testing verifier:", err)
		return false
	}

	responseBody := types.BinaryResponse{}
	responseDecoder := json.NewDecoder(response.Body)
	if err := responseDecoder.Decode(&responseBody); err != nil {
		fmt.Println("Error decoding response while testing verifier:", err)
		return false
	}

	return responseBody.Result
}

func (g *RemoteGame) MakeGuess(player Player, code []int) bool {
	request := types.MakeGuessRequest{
		GameIndex: g.gameIndex,
		Code:      code,
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling request while making guess:", err)
		return false
	}

	response, err := g.client.Post(g.addr+"/player/make-guess", "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		fmt.Println("error making guess:", err)
		return false
	} else if err := checkResponse(response); err != nil {
		fmt.Println("Response error making guess:", err)
		return false
	}

	responseBody := types.BinaryResponse{}
	responseDecoder := json.NewDecoder(response.Body)
	if err := responseDecoder.Decode(&responseBody); err != nil {
		fmt.Println("Error decoding response while making guess:", err)
		return false
	}

	return responseBody.Result
}

func (g *RemoteGame) Rank() [][]Player {
	return nil
}

func (g *RemoteGame) Stats() map[Player]*PlayerMoves {
	return nil
}
