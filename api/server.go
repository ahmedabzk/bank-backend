package api

import (
	"fmt"

	db "github.com/ahmedabzk/simple_bank/db/sqlc"
	"github.com/ahmedabzk/simple_bank/token"
	"github.com/ahmedabzk/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service
type Server struct {
	config     util.Config
	store      *db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewServer creates new HTTP server and setup routing
func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TOKEN_SYMMETRIC_KEY)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/account", server.createAccount)
	authRoutes.GET("/account", server.listAccounts)
	authRoutes.GET("/account/:id", server.getAccount)

	authRoutes.POST("/transfer", server.createTransfer)

	server.router = router
}

// Start runs the server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func ErrorResponds(err error) gin.H {
	return gin.H{"err": err.Error()}
}
