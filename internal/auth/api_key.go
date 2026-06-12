package auth 

import (
	"net/http"
	"strings"
	"errors"
)

func GetAPIKey(headers http.Header) (string, error){

	apiString, ok := headers["Authorization"]
	if ok{
		return strings.Split(apiString[0], " ")[1], nil
	} else{
		
		return "", errors.New("no authorization header found")
	}
}