-- name: GetOneChirp :one

SELECT * from chirps 
WHERE id = $1; 