package main 

import (
	"net/http"
	"log"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/arunima1319/chirpy/internal/database"
	"github.com/arunima1319/chirpy/internal/auth"
)

type chirp struct{ 
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerDeleteChirp (w http.ResponseWriter, r *http.Request){

	token, err := auth.GetBearerToken(r.Header)
	if err!=nil{
		respondWithError(w, 401, "no token found in request header")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err!=nil{
		respondWithError(w, 401, "token could not be validated")
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err!=nil{
		log.Printf("chirp ID could not be parsed: %s", err)
		return
	}

	associatedUser, err := cfg.dbQueries.GetUserFromChirp(r.Context(), chirpID)
	if err!= nil{
		log.Printf("error in finding associated user: %s", err)
		respondWithError(w, 404, "the chirp is not found")
		return
	}

	if userID != associatedUser.ID{
		respondWithError(w, 403, "the user is not authorized to delete this chirp")
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), chirpID)
	if err!=nil{
		log.Printf("chirp could not be deleted: %s", err)
	}

	respondWithJSON(w, 204, map[string]string{"Success": "the chirp is deleted"})

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
	

	token, err := auth.GetBearerToken(r.Header)
	if err!=nil{
		log.Printf("error in getting bearer token: %s", err)
		return
	}
	userIDFromJWT, err := auth.ValidateJWT(token, cfg.secret)
	if err!=nil{
		respondWithError(w, 401, "JWT could not be validated")
		return
	}

	chirpResource, err := cfg.dbQueries.CreateChirp(
		r.Context(),
		database.CreateChirpParams{
			Body: cleanedText,
			UserID: userIDFromJWT, 
		})
	if err!=nil{
		log.Printf("Error in storing chirp to database: %s", err)
		return
	}
	
	chirpedResponse := chirp{
		ID: chirpResource.ID,
		CreatedAt: chirpResource.CreatedAt,
		UpdatedAt: chirpResource.UpdatedAt,
		UserID: userIDFromJWT,
		Body: cleanedText,

	}


	respondWithJSON(w, 201, chirpedResponse)
}