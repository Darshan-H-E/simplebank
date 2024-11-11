package api

import (
	"github.com/gin-gonic/gin"
	db "simplebank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

// new server instance
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// routes
	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccount)
  router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

// start server with address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
