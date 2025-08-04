package db

import (
	"context"
	"database/sql"
	"shopping/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomProduct(t *testing.T) Product {
	user := createRandomUser(t)

	arg := CreateProductParams{
		Name: util.RandomName(),
		Size: util.RandomProductSize(),
		Weight: util.RandomProductWeight(),
		Price: util.RandomProductPrice(),
		UserID: user.ID,
		CategoryID: util.RandomCategoryID(), // assign categoryID randomly
	}

	product, err := testQueries.CreateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product)
	
	// check for the params
	require.Equal(t, arg.Name, product.Name)
	require.Equal(t, arg.Size, product.Size)
	require.Equal(t, arg.Weight, product.Weight)
	require.Equal(t, arg.Price, product.Price)

	// check for timestamp and auto-generated product id
	require.NotZero(t, product.ID)
	require.NotZero(t, product.CreatedAt)
	require.NotZero(t, product.UpdatedAt)

	return product
}

func TestCreateProduct(t *testing.T) {
	createRandomProduct(t)
}

func TestGetProduct(t *testing.T) {
	product1 := createRandomProduct(t)
	product2, err := testQueries.GetProduct(context.Background(), product1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, product2)

	require.Equal(t, product1.ID, product2.ID)
	require.Equal(t, product1.Name, product2.Name)
	require.Equal(t, product1.Size, product2.Size)
	require.Equal(t, product1.Weight, product2.Weight)
	require.Equal(t, product1.Price, product2.Price)
	require.WithinDuration(t, product1.CreatedAt, product2.CreatedAt, time.Second)
	require.WithinDuration(t, product1.UpdatedAt, product2.UpdatedAt, time.Second)
}

func TestUpdateProduct(t *testing.T) {
	product1 := createRandomProduct(t)

	arg := UpdateProductParams{
		ID: product1.ID,
		Name: util.RandomName(),
		Size: util.RandomProductSize(),
		Weight: util.RandomProductWeight(),
		Price: util.RandomProductPrice(),
	}

	product2, err := testQueries.UpdateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, product2)

	require.Equal(t, product1.ID, product2.ID)
	require.Equal(t, arg.Name, product2.Name)
	require.Equal(t, arg.Size, product2.Size)
	require.Equal(t, arg.Weight, product2.Weight)
	require.Equal(t, arg.Price, product2.Price)
	require.WithinDuration(t, product1.CreatedAt, product2.CreatedAt, time.Second)
	require.WithinDuration(t, product1.UpdatedAt, product2.UpdatedAt, time.Second)
}

func TestDeleteProduct(t *testing.T) {
	product1 := createRandomProduct(t)
	arg := DeleteProductParams{
		ID: product1.ID,
	}
	err := testQueries.DeleteProduct(context.Background(), arg)
	require.NoError(t, err)

	product2, err := testQueries.GetProduct(context.Background(), product1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, product2)
}

func TestListProduct(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomProduct(t)
	}

	arg := ListProductsParams{
		Limit: 5,
		Offset: 5,
	}

	products, err := testQueries.ListProducts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, products, 5)

	for _, product := range products {
		require.NotEmpty(t, product)
	}
}