package api

import (
	db "github.com/ahmedabzk/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.POST("/account", server.CreateAccount)
	server.router = router
	return server
}

// Start runs the server on a specific address
func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

func ErrorResponds(err error) gin.H{
	return gin.H{"err": err.Error()}
}
