package main 

import (
	"net/http"
	"log"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/arunima1319/chirpy/internal/database"
)

type chirp struct{ 
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`

}

func (cfg *apiConfig) handlerGetChirps (w http.ResponseWriter, r *http.Request){
	chirpSlice, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err!=nil{
		log.Printf("Error in getting chirps: %s", err)
	} 



	chirpList := []chirp{}

	for _, singleChirp := range chirpSlice{
		chirpJSON := chirp{
			ID: singleChirp.ID,
			CreatedAt: singleChirp.CreatedAt, 
			UpdatedAt: singleChirp.UpdatedAt, 
			Body: singleChirp.Body,
			UserID: singleChirp.UserID,
		}
		
		chirpList = append(chirpList, chirpJSON)
	}

	respondWithJSON(w, 200, chirpList)
}

func (cfg *apiConfig) handlerGetOneChirp (w http.ResponseWriter, r *http.Request){
	uuidID, err := uuid.Parse(r.PathValue("chirpID"))
	if err!= nil{
		log.Printf("error in parsing the ID string: %s", err)
	}
	chirpDatabase, err := cfg.dbQueries.GetOneChirp(r.Context(), uuidID)
	if err!=nil{
		log.Printf("error in getting the chirp from database: %s", err)
		respondWithError(w, 404, "Resource not found")
		return
	}

	chirpJSON := chirp{
		ID: chirpDatabase.ID, 
		CreatedAt: chirpDatabase.CreatedAt, 
		UpdatedAt: chirpDatabase.UpdatedAt, 
		Body: chirpDatabase.Body, 
		UserID: chirpDatabase.UserID,
	}

	respondWithJSON(w, 200, chirpJSON)
}

func (cfg *apiConfig) handlerChirps (w http.ResponseWriter, r *http.Request){ 

	chirped := chirp{}

	decoder := json.NewDecoder(r.Body)
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

	cleanedText := replaceBadWord(chirped.Body)
	chirped.Body = cleanedText

	chirpResource, err := cfg.dbQueries.CreateChirp(
		r.Context(),
		database.CreateChirpParams{
			Body: cleanedText,
			UserID: chirped.UserID, 
		})
	if err!=nil{
		log.Printf("Error in storing chirp to database: %s", err)
		return
	}

	chirped.ID = chirpResource.ID
	chirped.CreatedAt = chirpResource.CreatedAt
	chirped.UpdatedAt = chirpResource.UpdatedAt

	respondWithJSON(w, 201, chirped)
}