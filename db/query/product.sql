-- name: CreateProduct :one
INSERT INTO products (
  name,
  size,
  weight,
  price,
  user_id,
  category_id
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;


-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListProductsWithCategory :many
SELECT * FROM products
JOIN categories ON products.category_id = categories.id
WHERE categories.name = $1
ORDER BY products.id
LIMIT $2
OFFSET $3;

-- name: UpdateProduct :one
UPDATE products
  set name = $2,
  size = $3,
  weight = $4,
  price = $5
WHERE id = $1
AND user_id = $6
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1 AND user_id = $2;