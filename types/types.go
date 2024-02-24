package types

type JoinRequest struct {
	PlayerName string `json:"playerName"`
}

type GetGamesResponse struct {
	Games [][]int `json:"games"`
}

type AskQuestionRequest struct {
	GameIndex     int   `json:"gameIndex"`
	VerifierIndex int   `json:"verifierIndex"`
	Code          []int `json:"code"`
}

type MakeGuessRequest struct {
	GameIndex int   `json:"gameIndex"`
	Code      []int `json:"code"`
}

type BinaryResponse struct {
	Result bool `json:"result"`
}

type RankRequest struct {
	GameIndex int `json:"gameIndex"`
}

type RankResponse struct {
	Rankings [][]*RemotePlayer `json:"rankings"`
}

type RemotePlayer struct {
	Name string `json:"name"`
}

func (p *RemotePlayer) GetPlayerName() string {
	return p.Name
}
