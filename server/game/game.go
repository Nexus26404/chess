package game

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	BoardSize = 15 // Standard 15x15 board
)

var gameStore = make(map[string]*Game)

type Game struct {
	Room          *Room
	ID            string
	Board         [][]int
	CurrentPlayer PlayerTurn
	Winner        PlayerTurn
	IsFinished    bool
}

func NewGame() *Game {
	board := make([][]int, BoardSize)
	for i := range board {
		board[i] = make([]int, BoardSize)
	}

	game := &Game{
		ID:            generateGameID(),
		Board:         board,
		CurrentPlayer: PlayerTurnBlack,
		Winner:        PlayerTurnNone,
		IsFinished:    false,
	}
	gameStore[game.ID] = game
	return game
}

func (g *Game) MakeMove(row, col int) bool {
	if row < 0 || row >= BoardSize || col < 0 || col >= BoardSize || g.Board[row][col] != int(PlayerTurnNone) || g.IsFinished {
		return false
	}

	g.Board[row][col] = int(g.CurrentPlayer)
	if g.checkWin(row, col) {
		g.Winner = g.CurrentPlayer
		g.IsFinished = true
	} else {
		g.CurrentPlayer = 3 - g.CurrentPlayer // Switch between 1 and 2
	}
	gameStore[g.ID] = g
	return true
}

func (g *Game) checkWin(row, col int) bool {
	directions := [][2]int{{1, 0}, {0, 1}, {1, 1}, {1, -1}}
	player := g.Board[row][col]

	for _, dir := range directions {
		count := 1
		// Check forward
		for i := 1; i < 5; i++ {
			r, c := row+dir[0]*i, col+dir[1]*i
			if !isValidPos(r, c) || g.Board[r][c] != player {
				break
			}
			count++
		}
		// Check backward
		for i := 1; i < 5; i++ {
			r, c := row-dir[0]*i, col-dir[1]*i
			if !isValidPos(r, c) || g.Board[r][c] != player {
				break
			}
			count++
		}
		if count >= 5 {
			return true
		}
	}
	return false
}

func (g *Game) SwitchTurn() {
	g.Room.mutex.Lock()
	defer g.Room.mutex.Unlock()
	if g.CurrentPlayer == PlayerTurnBlack {
		g.CurrentPlayer = PlayerTurnWhite
	} else {
		g.CurrentPlayer = PlayerTurnBlack
	}
}

func (r *Room) CanMove(playerID string) bool {
	for _, player := range r.Players {
		if player.UserID == playerID {
			return r.Game.CurrentPlayer == player.Turn
		}
	}
	return false
}

func isValidPos(row, col int) bool {
	return row >= 0 && row < BoardSize && col >= 0 && col < BoardSize
}

func generateGameID() string {
	return "game-" + uuid.New().String()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust this for production
	},
}

type GameMessage struct {
	Type       string     `json:"type"`
	GameID     string     `json:"gameId"`
	Row        int        `json:"row"`
	Col        int        `json:"col"`
	Board      [][]int    `json:"board"`
	Player     int        `json:"player"`
	Winner     int        `json:"winner"`
	RoomStatus RoomStatus `json:"roomStatus"`
	PlayerTurn PlayerTurn `json:"playerTurn"`
	Success    bool       `json:"success"`
	Error      string     `json:"error,omitempty"`
	Owner      string     `json:"owner"`
	RoomName   string     `json:"roomName"` // Add this field
}

func (g *Game) Reset() {
	g.Board = make([][]int, BoardSize)
	for i := range g.Board {
		g.Board[i] = make([]int, BoardSize)
	}
	g.CurrentPlayer = PlayerTurnBlack
	g.Winner = PlayerTurnNone
	g.IsFinished = false
	gameStore[g.ID] = g
}

func (r *Room) GetRoomUpdateMessage() GameMessage {
	return GameMessage{
		Type:       "room_update",
		RoomStatus: r.Status,
		PlayerTurn: r.Game.CurrentPlayer,
		Owner:      r.Owner(),
	}
}

func (g *Game) GetUpdateMessage() GameMessage {
	return GameMessage{
		Type:   "update",
		GameID: g.ID,
		Board:  g.Board,
		Player: int(g.CurrentPlayer),
		Winner: int(g.Winner),
	}
}

func HandleGameWebSocket(c *gin.Context) {
	roomID := c.Param("roomId")
	userID := c.Query("userId")
	log.Printf("WebSocket connection attempt - Room: %s, User: %s", roomID, userID)

	room, exists := GetRoom(roomID)
	if !exists {
		log.Printf("WebSocket connection failed - Room not found: %s", roomID)
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed - Room: %s, User: %s, Error: %s", roomID, userID, err)
		return
	}

	room.AddConnection(userID, ws)
	log.Printf("WebSocket connected - Room: %s, User: %s", roomID, userID)
	defer func() {
		ws.Close()
		room.RemoveConnection(userID)
		log.Printf("WebSocket disconnected - Room: %s, User: %s", roomID, userID)
	}()

	// Send initial room state
	ws.WriteJSON(room.GetRoomUpdateMessage())

	for {
		var msg GameMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error - Room: %s, User: %s, Error: %s", roomID, userID, err)
			break
		}

		switch msg.Type {
		case "ready":
			log.Printf("Player ready - Room: %s, User: %s", roomID, userID)
			room.SetReady(userID)
			room.BroadcastMessage(room.GetRoomUpdateMessage())
		case "cancel_ready":
			room.CancelReady(userID)
			room.BroadcastMessage(room.GetRoomUpdateMessage())
		case "start_game":
			log.Printf("Game start attempt - Room: %s, User: %s", roomID, userID)
			if !room.CanStart(userID) {
				log.Printf("Game start failed - Room: %s, User: %s - Not owner or players not ready", roomID, userID)
				ws.WriteJSON(GameMessage{
					Type:    "error",
					Error:   "Cannot start game - not owner or players not ready",
					Success: false,
				})
				continue
			}
			game := NewGame()
			room.Game = game
			room.SetGameStarted()
			log.Printf("Game started - Room: %s, Game: %s", roomID, game.ID)

			updateMsg := GameMessage{
				Type:       "game_start",
				GameID:     game.ID,
				Board:      game.Board,
				Player:     int(game.CurrentPlayer),
				RoomStatus: room.Status,
				PlayerTurn: game.CurrentPlayer,
				Owner:      room.Owner(),
				Success:    true,
			}
			room.BroadcastMessage(updateMsg)
		case "move":
			if room.Game == nil || !room.CanMove(userID) {
				log.Printf("Invalid move attempt - Room: %s, User: %s", roomID, userID)
				continue
			}
			if room.Game.MakeMove(msg.Row, msg.Col) {
				log.Printf("Move made - Room: %s, Game: %s, User: %s, Position: [%d,%d]",
					roomID, room.Game.ID, userID, msg.Row, msg.Col)
				room.Game.SwitchTurn()
				response := room.Game.GetUpdateMessage()
				response.PlayerTurn = room.Game.CurrentPlayer
				room.BroadcastMessage(response)
			}
		case "reset":
			log.Printf("Game reset - Room: %s, User: %s", roomID, userID)
			if room.Game != nil {
				room.Game.Reset()
				room.BroadcastMessage(room.Game.GetUpdateMessage())
			}
		case "get_state":
			if room.Game != nil {
				ws.WriteJSON(room.Game.GetUpdateMessage())
			}
			// Always send room state
			ws.WriteJSON(room.GetRoomUpdateMessage())
		}
	}
}

func CreateRoom(c *gin.Context) {
	userID := c.Query("userId")
	roomName := c.Query("roomName")
	room := NewRoom(userID)
	room.Name = roomName
	log.Printf("Room created: %s (%s) by user: %s", room.Name, room.ID, userID)
	c.JSON(http.StatusOK, gin.H{
		"roomId":   room.ID,
		"roomName": room.Name,
		"player":   PlayerTurnBlack,
	})
}

func JoinRoom(c *gin.Context) {
	userID := c.Query("userId")
	roomID := c.Param("roomId")
	log.Printf("Join room attempt - Room: %s, User: %s", roomID, userID)

	room, exists := GetRoom(roomID)
	if !exists {
		log.Printf("Join room failed - Room not found: %s", roomID)
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	if err := room.JoinRoom(userID); err != nil {
		log.Printf("Join room failed - Room: %s, User: %s, Error: %s", roomID, userID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("User %s successfully joined room %s", userID, roomID)
	c.JSON(http.StatusOK, gin.H{
		"roomId": room.ID,
		"player": PlayerTurnWhite,
	})
}

type RoomInfo struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Status RoomStatus `json:"status"`
	Owner  string     `json:"owner"`
	IsFull bool       `json:"isFull"`
}

func ListRooms(c *gin.Context) {
	roomMutex.RLock()
	rooms := make([]RoomInfo, 0, len(roomStore))
	for _, room := range roomStore {
		rooms = append(rooms, RoomInfo{
			ID:     room.ID,
			Name:   room.Name,
			Status: room.Status,
			Owner:  room.Owner(),
			IsFull: len(room.Players) >= room.MaxPlayers,
		})
	}
	roomMutex.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
	})
}
