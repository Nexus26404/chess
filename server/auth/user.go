package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Nickname    string `json:"nickname"`
	GamesPlayed int    `json:"gamesPlayed"`
	Wins        int    `json:"wins"`
	Losses      int    `json:"losses"`
	Draws       int    `json:"draws"`
}

var users []User

func init() {
	// Load users from the JSON file during initialization
	file, err := os.Open("users.json")
	if err != nil {
		panic("Failed to load users: " + err.Error())
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&users); err != nil {
		panic("Failed to parse users: " + err.Error())
	}
}

func AuthenticateUser(c *gin.Context) {
	var credentials User
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	for _, user := range users {
		if user.Username == credentials.Username && user.Password == credentials.Password {
			// Set session
			session := sessions.Default(c)
			session.Set("username", user.Username)
			session.Set("expires_at", time.Now().Add(time.Hour*24))

			if err := session.Save(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to save session",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
			return
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
}

func GetUserInfo(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username").(string)

	// Find user info (excluding password)
	for _, user := range users {
		if user.Username == username {
			c.JSON(http.StatusOK, gin.H{
				"id":          user.ID,
				"username":    user.Username,
				"nickname":    user.Nickname,
				"gamesPlayed": user.GamesPlayed,
				"wins":        user.Wins,
				"losses":      user.Losses,
				"draws":       user.Draws,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		expiresAt := session.Get("expires_at")

		if username == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check if session has expired
		if expiresAt != nil {
			if expiry, ok := expiresAt.(time.Time); ok {
				if time.Now().After(expiry) {
					session.Clear()
					session.Save()
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
