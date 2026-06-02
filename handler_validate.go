package main 

import (
	"net/http"
	"encoding/json"
	"log"
	"strings"
)




func handlerValidateChirp(w http.ResponseWriter, r *http.Request){ 

	type chirp struct{ 
		Body string `json:"body"`
	}

	type responseError struct{ 
		Error string `json:"error"`
	}
	type responseValid struct{
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirped := chirp{}
	err := decoder.Decode(&chirped)
	if err!= nil{ 
		log.Printf("Error in decoding JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(chirped.Body) > 140{ 
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedChirp := replaceBadWord(chirped.Body)

	respBody := responseValid{ 
		CleanedBody: cleanedChirp,
	}

	respondWithJSON(w, 200, respBody)
}


func replaceBadWord (s string) string{ 
	stringSlice := strings.Split(s, " ")
	cleanedStringSlice := []string{}
	for _, word := range stringSlice{ 
	
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax"{
			cleanedStringSlice = append(cleanedStringSlice, "****")
		}else{
			cleanedStringSlice = append(cleanedStringSlice, word)
		}
	}
	return strings.Join(cleanedStringSlice, " ")

}

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