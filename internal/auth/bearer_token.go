package auth 

import (
	"strings"
	"net/http"
	"fmt"
)

func GetBearerToken(headers http.Header) (string, error){ 
	authString, ok := headers["Authorization"]
	if ok{
		return strings.Split(authString[0], " ")[1], nil
	}else{
		return "", fmt.Errorf("no Authorization header found")
	}

	
}