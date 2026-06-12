package main 

import (
	"net/http"
	"encoding/json"
	"log"
	"errors"
	"database/sql"
	"github.com/arunima1319/chirpy/internal/auth"
	
)

func (cfg *apiConfig) handlerPolkaWebhooks (w http.ResponseWriter, r *http.Request){

	apiKey, err := auth.GetAPIKey(r.Header) 
	if err!=nil{
		log.Printf("error in getting api Key from request header: %s", err)
		respondWithError(w, 401, "Polka could not be authenticated")
		return
	}

	if apiKey != cfg.polkaKey{
		respondWithError(w, 401, "Polka could not be authenticated")
		return
	}

	webhookEvent := webhooksReq{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&webhookEvent)
	if err!=nil{
		log.Printf("error in decoding request body: %s", err)
	}

	if webhookEvent.Event != "user.upgraded"{
		respondWithJSON(w, 204, struct{}{})
	}else{
		_, err := cfg.dbQueries.UpgradeUserToChirpy(r.Context(), webhookEvent.Data.UserId)
		if err!= nil{
			log.Printf("error in upgrading user to chirpy: %s", err)
			if errors.Is(err, sql.ErrNoRows){
				respondWithError(w, 404, "user cannot be found")
			}
			return
		}

		respondWithJSON(w, 204, struct{}{})
	}




}