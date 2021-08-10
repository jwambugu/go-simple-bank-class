package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/jwambugu/go-simple-bank-class/db/sqlc"
	"github.com/jwambugu/go-simple-bank-class/token"
	"github.com/jwambugu/go-simple-bank-class/util"
)

// Server serves all HTTP requests for the banking service
type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func (server *Server) setupRouter() {
	router := gin.Default()

	v1 := router.Group("/v1")

	authRoutes := v1.Group("/").Use(authMiddleware(server.tokenMaker))
	auth := v1.Group("/auth")

	auth.POST("login", server.loginUser)
	authRoutes.GET("/accounts", server.getAccounts)
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccountByID)

	authRoutes.POST("/transfers", server.createTransfer)

	v1.POST("/users", server.createUser)

	server.router = router
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("currency", validCurrency); err != nil {
			return nil, fmt.Errorf("failed to register the currency validator: %v", err)
		}
	}

	server.setupRouter()
	return server, nil
}

// Start will run the http server on the specified address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
