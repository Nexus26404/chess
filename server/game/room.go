package game

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RoomStatus int
type PlayerTurn int

const (
	RoomStatusUnready RoomStatus = iota
	RoomStatusReady
	RoomStatusStarted
	RoomStatusFinished
)

const (
	PlayerTurnNone PlayerTurn = iota
	PlayerTurnBlack
	PlayerTurnWhite
)

type Player struct {
	UserID   string
	IsReady  bool
	IsOnline bool
}

type Room struct {
	ID          string
	Game        *Game
	Player1     *Player
	Player2     *Player
	Owner       string
	Status      RoomStatus
	PlayerTurn  PlayerTurn
	mutex       sync.Mutex
	connections map[string]*websocket.Conn // Map of userID to websocket connection
}

var (
	roomStore = make(map[string]*Room)
	roomMutex sync.RWMutex
)

func NewRoom(userId string) *Room {
	room := &Room{
		ID:   generateRoomID(),
		Game: NewGame(),
		Player1: &Player{
			UserID:   userId,
			IsReady:  false,
			IsOnline: true,
		},
		Owner:       userId,
		Status:      RoomStatusUnready,
		PlayerTurn:  PlayerTurnNone,
		connections: make(map[string]*websocket.Conn),
	}

	roomMutex.Lock()
	roomStore[room.ID] = room
	roomMutex.Unlock()
	return room
}

func (r *Room) JoinRoom(playerID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.Player2 != nil {
		return errors.New("room is full")
	}
	if r.Player1.UserID == playerID {
		return errors.New("player already in room")
	}

	r.Player2 = &Player{
		UserID:   playerID,
		IsReady:  false,
		IsOnline: true,
	}
	return nil
}

func (r *Room) SetReady(playerID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.Player2 != nil && r.Player2.UserID == playerID {
		r.Player2.IsReady = true
	}

	if r.Player2 != nil && r.Player2.IsReady {
		r.Status = RoomStatusReady
	}
}

func (r *Room) CancelReady(playerID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.Player2 != nil && r.Player2.UserID == playerID {
		r.Player2.IsReady = false
	}
	r.Status = RoomStatusUnready
}

func (r *Room) IsGameReady() bool {
	return r.Status == RoomStatusReady
}

func (r *Room) SetGameStarted() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.Status = RoomStatusStarted
	r.PlayerTurn = PlayerTurnBlack // Black moves first
}

func (r *Room) SetGameFinished() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.Status = RoomStatusFinished
}

func (r *Room) SwitchTurn() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.PlayerTurn == PlayerTurnBlack {
		r.PlayerTurn = PlayerTurnWhite
	} else {
		r.PlayerTurn = PlayerTurnBlack
	}
}

func (r *Room) CanMove(playerID string) bool {
	return (r.PlayerTurn == PlayerTurnBlack && playerID == r.Player1.UserID) ||
		(r.PlayerTurn == PlayerTurnWhite && r.Player2 != nil && playerID == r.Player2.UserID)
}

func (r *Room) GetPlayerColor(playerID string) int {
	if playerID == r.Player1.UserID {
		return BlackPlayer
	}
	if r.Player2 != nil && playerID == r.Player2.UserID {
		return WhitePlayer
	}
	return EmptyCell
}

func (r *Room) AddConnection(userID string, conn *websocket.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.connections[userID] = conn
}

func (r *Room) RemoveConnection(userID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.connections, userID)
}

func (r *Room) BroadcastMessage(message interface{}) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, conn := range r.connections {
		if conn != nil {
			conn.WriteJSON(message)
		}
	}
}

func (r *Room) IsOwner(playerID string) bool {
	return r.Owner == playerID
}

func (r *Room) CanStart(playerID string) bool {
	return r.IsOwner(playerID) && r.IsGameReady()
}

func GetRoom(roomID string) (*Room, bool) {
	roomMutex.RLock()
	defer roomMutex.RUnlock()
	room, exists := roomStore[roomID]
	return room, exists
}

func generateRoomID() string {
	return "room-" + uuid.New().String()
}
