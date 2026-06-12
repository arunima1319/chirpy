package main

import (
	"github.com/google/uuid"
	"time"
)


type userBody struct{
	Email string `json:"email"`
	Password string `json:"password"`
	//ExpirationTimeSeconds int `json:"expires_in_seconds"`
}

type User struct{
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	TokenString string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type chirp struct{ 
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}


