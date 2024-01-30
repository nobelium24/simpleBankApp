package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:ID", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/update", server.updateAccount)
	router.DELETE("/delete/:ID", server.deleteAccount)


	server.router = router
	return server
}

func (server* Server) Start (address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
