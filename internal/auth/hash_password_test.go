package auth 

import (
	"testing"
	"github.com/alexedwards/argon2id"
)

func TestHashPassword(t *testing.T){
	want, err := argon2id.CreateHash("randompassword", argon2id.DefaultParams)
	if err!=nil{
		t.Errorf("error in using the argon2id library")
	}
	have, err := HashPassword("randompassword")
	if err!= nil{
		t.Errorf("error in using the HashPassword function")
	}
	if want!=have {
		t.Errorf("it does not match")
	}
}