package main

import (
	"net/http"
	"log"
)

func (cfg *apiConfig) handlerResetDatabase (w http.ResponseWriter, r *http.Request){

	if cfg.platform != "dev"{
		respondWithError(w, 403, "Forbidden")
	}
	err := cfg.dbQueries.ResetDatabase(r.Context())
	if err!=nil{
		log.Printf("Error in resetting database: %s", err)
	}
}