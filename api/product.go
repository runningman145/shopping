package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "shopping/db/sqlc"
	"shopping/token"

	"github.com/gin-gonic/gin"
)

type createProductRequest struct {
	Name         string `json:"name" binding:"required"`
	Size         string `json:"size" binding:"required"`
	Weight       int64  `json:"weight"`
	Price        int64  `json:"price" binding:"required"`
	CategoryName string `json:"category_name" binding:"required"`
}

type createProductResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Size         string `json:"size"`
	Weight       int64  `json:"weight"`
	Price        int64  `json:"price"`
	CategoryName string `json:"category_name"`
}

func (server *Server) createProduct(ctx *gin.Context) {
	var req createProductRequest

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	category, err := server.store.GetCategoryByName(ctx, req.CategoryName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get user by username to get the user ID
	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateProductParams{
		Name:       req.Name,
		Size:       req.Size,
		Weight:     req.Weight,
		Price:      req.Price,
		UserID:     user.ID,
		CategoryID: category.ID,
	}

	// insert product into database if request is okay
	product, err := server.store.CreateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := createProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Size:         product.Size,
		Weight:       product.Weight,
		Price:        product.Price,
		CategoryName: category.Name,
	}

	ctx.JSON(http.StatusOK, res)
}

type getProductRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getProduct(ctx *gin.Context) {
	var req getProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get the authenticated user from the authorization payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	product, err := server.store.GetProduct(ctx, req.ID)
	if err != nil { // error querying from the database or product with specific id doesn't exist
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get user by username to get the user ID for authorization check
	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if the authenticated user owns this product
	if product.UserID != user.ID {
		err := errors.New("you can only view your own products")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// delete product by id
type deleteProductRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteProduct(ctx *gin.Context) {
	var req deleteProductRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get the authenticated user from the authorization payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// get the product to check ownership
	product, err := server.store.GetProduct(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get user by username to get the user ID for authorization check
	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if the authenticated user owns this product
	if product.UserID != user.ID {
		err := errors.New("you can only delete your own products")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	err = server.store.DeleteProduct(ctx, db.DeleteProductParams{ID: req.ID})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, "Product deleted!")
}

// updateProduct
type updateProductRequest struct {
	ID     int64  `json:"id" binding:"required"`
	Name   string `json:"name"`
	Size   string `json:"size"`
	Weight int64  `json:"weight"`
	Price  int64  `json:"price"`
}

func (server *Server) updateProduct(ctx *gin.Context) {
	var req updateProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get the authenticated user from the authorization payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// get the product to check ownership
	product, err := server.store.GetProduct(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get user by username to get the user ID for authorization check
	user, err := server.store.GetUser(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check if the authenticated user owns this product
	if product.UserID != user.ID {
		err := errors.New("you can only update your own products")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	arg := db.UpdateProductParams{
		ID:     req.ID,
		Name:   req.Name,
		Size:   req.Size,
		Weight: req.Weight,
		Price:  req.Price,
	}

	updatedProduct, err := server.store.UpdateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, updatedProduct)
}

// list all products and only return the name and price
type listProductRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type listProductResponse struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

func (server *Server) listProducts(ctx *gin.Context) {
	var req listProductRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListProductsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize, // number of records database should skip
	}

	products, err := server.store.ListProducts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var res []listProductResponse
	for _, product := range products {
		res = append(res, listProductResponse{
			Name:  product.Name,
			Price: product.Price,
		})
	}

	ctx.JSON(http.StatusOK, res)
}
