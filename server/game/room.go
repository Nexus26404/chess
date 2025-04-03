package game

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

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
	UserID  string
	Name    string
	IsReady bool
	IsOwner bool
	Turn    PlayerTurn
}

type Room struct {
	ID          string
	Name        string
	Game        *Game
	Status      RoomStatus
	mutex       sync.Mutex
	connections map[string]*websocket.Conn // Map of userID to websocket connection
	MaxPlayers  int
	Players     []*Player
}

func (r *Room) Owner() string {
	for _, v := range r.Players {
		if v.IsOwner {
			return v.UserID
		}
	}
	return ""
}

func (r *Room) AllReady() bool {
	for _, player := range r.Players {
		if !player.IsReady {
			return false
		}
	}
	return true
}

var (
	roomStore = make(map[string]*Room)
	roomMutex sync.RWMutex
)

func NewRoom(userId string) *Room {
	room := &Room{
		ID:          generateRoomID(),
		Name:        "", // Will be set by CreateRoom
		Game:        NewGame(),
		Status:      RoomStatusUnready,
		connections: make(map[string]*websocket.Conn),
		MaxPlayers:  2,
		Players:     make([]*Player, 0, 2),
	}

	room.Players = append(room.Players, &Player{
		UserID:  userId,
		IsOwner: true,
		IsReady: true, // Owner is ready by default
	})

	roomMutex.Lock()
	roomStore[room.ID] = room
	roomMutex.Unlock()
	return room
}

func (r *Room) JoinRoom(playerID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(r.Players) >= r.MaxPlayers {
		return errors.New("room is full")
	}

	// Check if the player is already in the room
	for _, player := range r.Players {
		if player.UserID == playerID {
			return errors.New("player already in room")
		}
	}

	// Add the player to the room
	r.Players = append(r.Players, &Player{
		UserID:  playerID,
		IsReady: false,
	})

	return nil
}

func (r *Room) SetReady(playerID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, player := range r.Players {
		if player.UserID == playerID {
			player.IsReady = true
			break
		}
	}

	if r.AllReady() {
		r.Status = RoomStatusReady
	}
}

func (r *Room) CancelReady(playerID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, player := range r.Players {
		if player.UserID == playerID {
			player.IsReady = false
			break
		}
	}

	if !r.AllReady() {
		r.Status = RoomStatusUnready
	}
}

func (r *Room) IsGameReady() bool {
	return r.Status >= RoomStatusReady
}

func (r *Room) SetGameStarted() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Create a local random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if rng.Float32() < 0.5 {
		// Randomly assign player colors
		r.Players[0].Turn = PlayerTurnBlack
		r.Players[1].Turn = PlayerTurnWhite
	} else {
		r.Players[0].Turn = PlayerTurnWhite
		r.Players[1].Turn = PlayerTurnBlack
	}

	r.Status = RoomStatusStarted
}

func (r *Room) SetGameFinished() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.Status = RoomStatusFinished
}

func (r *Room) AddConnection(userID string, conn *websocket.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.connections[userID] = conn
}

func (r *Room) RemoveConnection(userID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	log.Printf("Removing connection for user %s from room %s", userID, r.ID)
	delete(r.connections, userID)

	// Find player index
	playerIndex := -1
	for i, player := range r.Players {
		if player.UserID == userID {
			playerIndex = i
			break
		}
	}

	if playerIndex != -1 {
		wasOwner := r.Players[playerIndex].IsOwner
		// Remove the player using slice operations
		r.Players = append(r.Players[:playerIndex], r.Players[playerIndex+1:]...)

		if wasOwner {
			if len(r.Players) == 0 {
				log.Printf("Owner left empty room %s, destroying room", r.ID)
				r.DestroyRoom()
			} else {
				// Transfer ownership to the first remaining player
				r.Players[0].IsOwner = true
				r.Players[0].IsReady = true
				log.Printf("Transferring room %s ownership to user %s", r.ID, r.Players[0].UserID)
			}
		}

		r.Status = RoomStatusUnready
		log.Printf("Player %s removed from room %s, %d players remaining", userID, r.ID, len(r.Players))
	}

	// Broadcast updated room state to remaining players
	r.BroadcastMessage(r.GetRoomUpdateMessage())
}

func (r *Room) DestroyRoom() {
	roomMutex.Lock()
	defer roomMutex.Unlock()
	delete(roomStore, r.ID)
}

func (r *Room) TransferOwnership(newOwnerID string) {
	for _, player := range r.Players {
		// Remove owner status from all players
		player.IsOwner = (player.UserID == newOwnerID)
	}
	// Reset ready status when ownership changes
	r.Status = RoomStatusUnready
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
	for _, player := range r.Players {
		if player.UserID == playerID && player.IsOwner {
			return true
		}
	}
	return false
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
