package main 

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/arunima1319/chirpy/internal/auth"
	"github.com/arunima1319/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser (w http.ResponseWriter, r *http.Request){ 

	reqBody := userBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err!=nil{
		log.Printf("Error in decoding request: %s", err)
		return
	}

	hashedPassword, err := auth.HashPassword(reqBody.Password)
	if err != nil{
		log.Printf("Error in hashing password: %s", err)
	}

	user, err := cfg.dbQueries.CreateUser(
		r.Context(), 
		database.CreateUserParams{
			Email: reqBody.Email,
			HashedPassword: hashedPassword,
		})
	if err!= nil{
		log.Printf("Error in creating user: %s", err)
		return
	}
	userToRespondWith := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt, 
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, 201, userToRespondWith)

}