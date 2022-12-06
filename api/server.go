package api

import (
	db "github.com/ahmedabzk/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/account", server.createAccount)
	router.POST("/transfer", server.createTransfer)

	router.GET("/account", server.listAccounts)
	router.GET("/account/:id", server.getAccount)
	server.router = router
	return server
}

// Start runs the server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func ErrorResponds(err error) gin.H {
	return gin.H{"err": err.Error()}
}
