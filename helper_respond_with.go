package main 

import (
	"net/http"
	"encoding/json"
	"log"
)
func respondWithError(w http.ResponseWriter, code int, msg string){
	respondWithJSON(w, code, map[string]string{"Error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){ 
	dat, err := json.Marshal(payload)
		if err!= nil{ 
			log.Printf("Error in Marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}