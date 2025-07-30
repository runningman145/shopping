package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
	// rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// Random generates a random name
func RandomName() string {
	return RandomString(8)
}

// RandomProduct size
func RandomProductSize() string {
	sizes := []string{"S", "M", "L", "XL"}
	n := len(sizes)
	
	return sizes[rand.Intn(n)]
}

// randomProduct weight
func RandomProductWeight() int64 {
	return RandomInt(500, 2500)
}

// randomProduct Price
func RandomProductPrice() int64 {
	return RandomInt(5000, 250000)
}

// randomcategory assignment to a product, should be one of the categories in the categories table
func RandomCategoryID() int64 {
	// For demonstration, return a random category ID between 1 and 10
	return RandomInt(1, 10)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomString(6))
}