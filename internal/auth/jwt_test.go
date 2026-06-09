package auth

import (
	"testing"
	"github.com/google/uuid"
	"time"
	
)

func TestMakeAndValidateJWT(t * testing.T){
	uuidTest := uuid.New()
	expiration, err := time.ParseDuration("300s")
	if err!=nil{
		t.Errorf("error in creating duration: %s", err)
	}
	jwtString, err := MakeJWT(uuidTest, "my-secret-go-token", expiration)
	if err!=nil{
		t.Errorf("error in making JWT: %s", err)
	}

	uuidActual, err := ValidateJWT(jwtString, "my-secret-go-token")
	if err!= nil{
		t.Errorf("error in validating JWT: %s", err)
	}

	if uuidTest != uuidActual{
		t.Errorf("The uuids do not match")
	}
}