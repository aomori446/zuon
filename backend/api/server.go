package api

import (
	"log"
	"net/http"

	"github.com/aomori446/zuon/internal/auth"
	"github.com/aomori446/zuon/internal/unsplash"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	client *unsplash.Client
}

func NewServer(apiKey string) (*Server, error) {
	if apiKey == "" {
		log.Println("Warning: UNSPLASH_ACCESS_KEY not set")
	}

	client, err := unsplash.NewClient(apiKey)
	if err != nil {
		return nil, err
	}

	r := gin.Default()

	// Simple CORS middleware for localhost
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
		c.Next()
	})

	s := &Server{
		router: r,
		client: client,
	}

	s.routes()

	return s, nil
}

func (s *Server) routes() {
	unsplashHandler := NewUnsplashHandler(s.client)
	authHandler := NewAuthHandler()
	
	// Auth routes
	authGroup := s.router.Group("/api/v1/auth/github")
	{
		authGroup.GET("/login", authHandler.Login)
		authGroup.GET("/callback", authHandler.Callback)
		authGroup.GET("/poll", authHandler.Poll)
	}

	// Protected Unsplash routes
	s.router.GET("/search", authMiddleware(), unsplashHandler.Search)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
