package main 

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/arunima1319/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request){

	reqBody := userBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err!=nil {
		log.Printf("Error decoding request")
	}
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), reqBody.Email)
	if err!=nil {
		respondWithError(w, 401, "incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(reqBody.Password, user.HashedPassword)
	if err!=nil{
		log.Printf("Error checking password hash")
		return
	}
	if match != true{
		respondWithError(w, 401, "incorrect email or password")
	}else{
		userResourceCopy := User{
			ID: user.ID, 
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		}
		respondWithJSON(w, 200, userResourceCopy)


	}

}