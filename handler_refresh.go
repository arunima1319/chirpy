package main 

import (
	"net/http"
	"github.com/arunima1319/chirpy/internal/auth"
	"time"
	"log"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request){

	tokenString, err := auth.GetBearerToken(r.Header)
	if err!= nil{
		log.Printf("error in getting refresh token from request header: %s", err)
	}
	refreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), tokenString)
	if err!=nil{
		respondWithError(w, 401, "Refresh token does not exist")
	}else if refreshToken.RevokedAt.Valid == true{
		respondWithError(w, 401, "Refresh token has been revoked")
	}else if time.Now().UTC().After(refreshToken.ExpiresAt.UTC()){
		respondWithError(w, 401, "Refresh token has expired")
	}else{
		duration, _ := time.ParseDuration("1h")
		accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, duration)
		if err!= nil{
			log.Printf("error in creating a new access token: %s", err)
		}else{ 
			respondWithJSON(w, 200, map[string]string{"token": accessToken})
		}
	}
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request){

	tokenString, err := auth.GetBearerToken(r.Header)
	if err!= nil{
		log.Printf("error in getting refresh token from request header: %s", err)
	}

	err = cfg.dbQueries.RevokeToken(r.Context(), tokenString)
	if err!=nil{
		respondWithError(w, 401, "refresh token not found")
	}

	respondWithJSON(w, 204, map[string]string{"Success": "Refresh token has been revoked"})

}