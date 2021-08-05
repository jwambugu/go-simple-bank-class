package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/jwambugu/go-simple-bank-class/db/sqlc"
	"log"
)

// Server serves all HTTP requests for the banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("currency", validCurrency); err != nil {
			log.Fatalf("failed to register the currency validator: %v", err)
		}
	}

	v1 := router.Group("/v1")
	accounts := v1.Group("/accounts")
	transfers := v1.Group("/transfers")
	users := v1.Group("/users")

	accounts.GET("", server.getAccounts)
	accounts.POST("", server.createAccount)
	accounts.GET(":id", server.getAccountByID)

	transfers.POST("", server.createTransfer)

	users.POST("", server.createUser)

	server.router = router
	return server
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
