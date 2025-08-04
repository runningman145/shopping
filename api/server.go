package api

import (
	"fmt"
	db "shopping/db/sqlc"
	"shopping/token"
	"shopping/util"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// server serves HTTP requests for our shopping service
type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      *db.Store
	router     *gin.Engine
}

func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.Token)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// Configure CORS
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

	router.POST("/users", server.createUser)
	router.POST("users/login", server.loginUser)

	// list categories, description, then all products in that category, list by name

	router.GET("/categories", server.listCategories)
	// router.GET("categories/:category/products", server.listCategoryProducts)
	router.GET("/products", server.listProducts)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/me", server.getCurrentUser)
	authRoutes.GET("/admin/users", server.listUsers)

	authRoutes.POST("/categories", server.createCategory)
	authRoutes.DELETE("/categories/:id", server.deleteCategory)
	authRoutes.GET("/categories/:id", server.getCategory)
	authRoutes.PUT("/categories/:id", server.updateCategory)

	authRoutes.POST("/products", server.createProduct)
	authRoutes.GET("/products/:id", server.getProduct)
	authRoutes.DELETE("/products/:id", server.deleteProduct)
	authRoutes.PUT("/products/:id", server.updateProduct)

	server.router = router
}

// Start runs the HTTP server on a specific address
// router is private: can't be accessed outside api
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
