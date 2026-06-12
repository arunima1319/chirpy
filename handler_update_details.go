package main 

import (
	"net/http"
	"log"
	"github.com/arunima1319/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpdateDetails (w http.ResponseWriter, r *http.Request){
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err!= nil{
		log.Printf("Error in getting token from request header: %s", err)
	}

	userID, err := auth.ValidateJWT(jwtToken, cfg.secret)
	if err!=nil{
		log.Printf("error in validating JWT: %s", err)
	}

	
}