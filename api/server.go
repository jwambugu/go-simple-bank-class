package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/jwambugu/go-simple-bank-class/db/sqlc"
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

	v1 := router.Group("/v1")

	v1.POST("/accounts", server.createAccount)
	v1.GET("/accounts", server.getAccounts)
	v1.GET("/accounts/:id", server.getAccountByID)

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
