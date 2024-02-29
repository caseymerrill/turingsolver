package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/caseymerrill/turingsolver/debounce"
	"github.com/caseymerrill/turingsolver/game"
	"github.com/caseymerrill/turingsolver/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type GameServer struct {
	games         []game.Game
	players       map[string]game.Player
	playersLock   sync.RWMutex
	printWinCount func()
}

func (s *GameServer) Join(c *gin.Context) {
	request := types.JoinRequest{}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error})
		fmt.Println(err)
		return
	}

	if request.PlayerName == "" {
		c.JSON(400, gin.H{"error": "Player name is required"})
		return
	}

	s.playersLock.Lock()
	defer s.playersLock.Unlock()
	if _, exists := s.players[request.PlayerName]; exists {
		c.JSON(400, gin.H{"error": "Player already exists"})
		return
	}

	s.players[request.PlayerName] = &types.RemotePlayer{Name: request.PlayerName}

	session := sessions.Default(c)
	session.Set("playerName", request.PlayerName)
	if err := session.Save(); err != nil {
		delete(s.players, request.PlayerName)
		c.JSON(500, gin.H{"error": "Failed to save session"})
	}

	c.JSON(200, gin.H{"playerName": request.PlayerName})
}

func (s *GameServer) GetGames(c *gin.Context) {
	response := types.GetGamesResponse{
		Games: make([][]int, len(s.games)),
	}

	for gameIndex := range s.games {
		cards := s.games[gameIndex].GetVerifierCards()
		response.Games[gameIndex] = make([]int, len(cards))
		for cardIndex, card := range cards {
			response.Games[gameIndex][cardIndex] = card.CardNumber
		}
	}

	c.JSON(200, response)
}

func (s *GameServer) AskQuestion(c *gin.Context) {
	request := types.AskQuestionRequest{}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error})
		fmt.Println(err)
		return
	} else if request.GameIndex < 0 || request.GameIndex >= len(s.games) {
		c.JSON(400, gin.H{"error": "Invalid game index"})
		return
	}

	currentGame := s.games[request.GameIndex]
	if request.VerifierIndex < 0 || request.VerifierIndex >= len(currentGame.GetVerifierCards()) {
		c.JSON(400, gin.H{"error": "Invalid verifier index"})
		return
	}

	player, ok := c.Get("player")
	if !ok {
		c.JSON(500, gin.H{"error": "Player not found"})
		return
	}

	check := currentGame.AskQuestion(player.(*types.RemotePlayer), request.Code, request.VerifierIndex)
	c.JSON(200, types.BinaryResponse{Result: check})
}

func (s *GameServer) MakeGuess(c *gin.Context) {
	request := types.MakeGuessRequest{}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error})
		fmt.Println(err)
		return
	} else if request.GameIndex < 0 || request.GameIndex >= len(s.games) {
		c.JSON(400, gin.H{"error": "Invalid game index"})
		return
	}

	player, ok := c.Get("player")
	if !ok {
		c.JSON(500, gin.H{"error": "Player not found"})
		return
	}

	currentGame := s.games[request.GameIndex]
	result := currentGame.MakeGuess(player.(*types.RemotePlayer), request.Code)

	s.printWinCount()

	c.JSON(200, types.BinaryResponse{Result: result})
}

func (s *GameServer) GetRank(c *gin.Context) {
	request := types.RankRequest{}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error})
		fmt.Println(err)
	} else if request.GameIndex < 0 || request.GameIndex >= len(s.games) {
		c.JSON(400, gin.H{"error": "Invalid game index"})
		return
	}

	currentGame := s.games[request.GameIndex]
	rankings := currentGame.Rank()

	// Convert to remote players
	remoteRankings := make([][]*types.RemotePlayer, len(rankings))
	for i, ranking := range rankings {
		remoteRankings[i] = make([]*types.RemotePlayer, len(ranking))
		for j, player := range ranking {
			remoteRankings[i][j] = player.(*types.RemotePlayer)
		}
	}

	c.JSON(200, types.RankResponse{Rankings: remoteRankings})
}

func (s *GameServer) Authenticate(c *gin.Context) {
	session := sessions.Default(c)
	playerName := session.Get("playerName")
	s.playersLock.RLock()
	defer s.playersLock.RUnlock()
	playerNameStr, ok := playerName.(string)
	if !ok {
		c.JSON(401, gin.H{"error": "Not authenticated, no player name set."})
		c.Abort()
		return
	}

	if player := s.players[playerNameStr]; player == nil {
		c.JSON(401, gin.H{"error": "Not authenticated, player not fount."})
		c.Abort()
		return
	} else {
		c.Set("player", player)
	}

	c.Next()
}

func (s *GameServer) Listen() {
	r := gin.Default()
	store := cookie.NewStore([]byte("super-secret-turing-game-cookie-key"))
	r.Use(sessions.Sessions("session", store))

	r.POST("/join", s.Join)

	authenticatedGroup := r.Group("/player", s.Authenticate)
	authenticatedGroup.GET("/games", s.GetGames)
	authenticatedGroup.POST("/test-verifier", s.AskQuestion)
	authenticatedGroup.POST("/make-guess", s.MakeGuess)
	authenticatedGroup.GET("/rank", s.GetRank)

	if err := r.Run(); err != nil {
		fmt.Println("Running server : ", err)
	}
}

func NewGameServer(games []game.Game) *GameServer {
	return &GameServer{
		games:       games,
		players:     make(map[string]game.Player),
		playersLock: sync.RWMutex{},
		printWinCount: debounce.Debounce(func() {
			game.PrintWinCount(games)
		}, 1*time.Second),
	}
}
