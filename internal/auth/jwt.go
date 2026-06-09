package auth 

import (
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func MakeJWT (userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){ 
	token := jwt.NewWithClaims(
				jwt.SigningMethodHS256, 
				jwt.RegisteredClaims{
					Issuer: "chripy-access",
					IssuedAt: &jwt.NumericDate{
						time.Now().UTC(),
					},
					ExpiresAt: &jwt.NumericDate{
						time.Now().Add(expiresIn).UTC(), 
					},
					Subject: userID.String(),
				})
	jwt, err := token.SignedString([]byte(tokenSecret))
	if err!= nil{
		return "", err
	}else{
		return jwt, nil
	}
}

func ValidateJWT (tokenString, tokenSecret string) (uuid.UUID, error){ 

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
						tokenString,
						claims,
						func (token *jwt.Token) (interface{}, error){
							return []byte(tokenSecret), nil
						})
	if err!=nil{
		return uuid.Nil, err
	}
	
	stringId, err := token.Claims.GetSubject()
	if err!=nil{
		return uuid.Nil, err
	}
	id, err := uuid.Parse(stringId)
	if err!=nil{
		return uuid.Nil, err
	}
	return id, nil
}
