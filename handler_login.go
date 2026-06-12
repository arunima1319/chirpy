package main 

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/arunima1319/chirpy/internal/auth"
	"github.com/arunima1319/chirpy/internal/database"
	"time"
	"database/sql"
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

	/*
	expirationDuration := "1h"
	ptrExpirationDuration := &expirationDuration
	
	if reqBody.ExpirationTimeSeconds == 0 || reqBody.ExpirationTimeSeconds > 3600{
		*ptrExpirationDuration = "1h"
	}else{
		*ptrExpirationDuration = fmt.Sprintf("%ds", reqBody.ExpirationTimeSeconds)
	}
	*/

	duration, err := time.ParseDuration("1h")
	if err!=nil{
		log.Printf("error in converting to a valid Duration")
	}

	
	match, err := auth.CheckPasswordHash(reqBody.Password, user.HashedPassword)
	if err!=nil{
		log.Printf("Error checking password hash")
		return
	}
	if match != true{
		respondWithError(w, 401, "incorrect email or password")
	}else{
		token, err := auth.MakeJWT(user.ID, cfg.secret, duration)
		if err!= nil{
			log.Printf("error in making JWT")
		}

		refreshToken, err := cfg.dbQueries.CreateRefreshToken(
								r.Context(), 
								database.CreateRefreshTokenParams{
									Token: auth.MakeRefreshToken(),
									CreatedAt: time.Now().UTC(),
									UpdatedAt: time.Now().UTC(), 
									UserID: user.ID, 
									ExpiresAt: time.Now().Add(time.Hour * 24 * 60).UTC(), 
									RevokedAt: sql.NullTime{Valid: false},
								})
		if err!= nil{
			log.Printf("Error in creating refresh token: %s", err)
		}
		userResourceCopy := User{
			ID: user.ID, 
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
			TokenString: token,
			RefreshToken: refreshToken.Token,
			IsChirpyRed: user.IsChirpyRed,

		}
		respondWithJSON(w, 200, userResourceCopy)


	}

}