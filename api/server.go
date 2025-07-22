package api

import (
	db "shopping/db/sqlc"

	"github.com/gin-gonic/gin"
)

// server serves HTTP requests for our shopping service
type Server struct {
	store *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// add routes to router
	router.POST("/products", server.createProduct)
	router.GET("/products/:id", server.getProduct)
	router.GET("/products", server.listProducts)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address
// router is private: can't be accessed outside api
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}