-- name: CreateUser :one

INSERT INTO users(id, created_at, updated_at, email)
VALUES(
    GEN_RANDOM_UUID(), 
    TIMESTAMP.NOW(), 
    TIMESTAMP.NOW(), 
    $1
)
RETURNING *;