-- name: GetUserFromChirp :one

SELECT users.* FROM users
INNER JOIN chirps  
    ON users.id = chirps.user_id
WHERE chirps.id = $1; 

