package api

import (
	"log"

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
	
	s.router.GET("/search", unsplashHandler.Search)
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
