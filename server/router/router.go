package router

import (
	"encoding/gob"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"chess/server/auth"
	"chess/server/game"
)

func RegisterRoutes(r *gin.Engine) {
	// Register time.Time type for session storage
	gob.Register(time.Time{})

	// Initialize session store
	store := cookie.NewStore([]byte("your_secret_key")) // Replace with a secure key
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Enable secure cookies
	})
	r.Use(sessions.Sessions("chess_session", store))

	// API routes
	r.POST("/authenticate", auth.AuthenticateUser)
	r.GET("/user/info", auth.AuthMiddleware(), auth.GetUserInfo)

	// Room routes
	r.GET("/room/list", auth.AuthMiddleware(), game.ListRooms)
	r.GET("/room/create", auth.AuthMiddleware(), game.CreateRoom)
	r.GET("/room/join/:roomId", auth.AuthMiddleware(), game.JoinRoom)
	r.GET("/room/ws/:roomId", auth.AuthMiddleware(), game.HandleGameWebSocket)

	// Serve static files from the Vue app's dist directory
	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")

	// Handle SPA routing - serve index.html for all other routes
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
}
