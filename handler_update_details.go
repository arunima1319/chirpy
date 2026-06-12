package main 

import (
	"net/http"
	"log"
	"github.com/arunima1319/chirpy/internal/auth"
	"encoding/json"
	"github.com/arunima1319/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateDetails (w http.ResponseWriter, r *http.Request){
	jwtToken, err := auth.GetBearerToken(r.Header)
	if err!= nil{
		log.Printf("Error in getting token from request header: %s", err)
	}

	userID, err := auth.ValidateJWT(jwtToken, cfg.secret)
	if err!=nil{
		respondWithError(w, 401, "JWT could not be validated")
		return
	}

	reqBody := userBody{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err!= nil{
		log.Printf("Error in decoding request body: %s", err)
		return
	}

	newEmail := reqBody.Email
	hashedPassword, err := auth.HashPassword(reqBody.Password)
	if err!= nil{
		log.Printf("Error in hashing password: %s", err)
		return
	}

	user, err := cfg.dbQueries.UpdateEmailAndPassword(
		r.Context(), 
		database.UpdateEmailAndPasswordParams{
			Email: newEmail,
			HashedPassword: hashedPassword,
			ID: userID,
		})
	if err!=nil{
		log.Printf("error in storing new email and password to database: %s", err)
		return
	}

	refreshToken, err := cfg.dbQueries.GetRefreshTokenFromUser(r.Context(), user.ID)
	if err!= nil{
		log.Printf("error in getting refresh token of user: %s", err)
	}
	userResource := User{
		ID: user.ID, 
		CreatedAt: user.CreatedAt, 
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		TokenString: jwtToken, 
		RefreshToken: refreshToken.Token,
	}
	respondWithJSON(w, 200, userResource)


}