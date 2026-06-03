package main 

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/google/uuid"
	"time"
)


type email struct{
	Email string `json:"email"`
}

type User struct{
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`

}

func (cfg *apiConfig) handlerCreateUser (w http.ResponseWriter, r *http.Request){ 

	emailAddress := email{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&emailAddress)
	if err!=nil{
		log.Printf("Error in decoding request: %s", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), emailAddress.Email)
	if err!= nil{
		log.Printf("Error in creating user: %s", err)
		return
	}
	userToRespondWith := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt, 
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, 201, userToRespondWith)

}