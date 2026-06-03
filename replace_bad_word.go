package main 

import "strings"


func replaceBadWord (s string) string{ 
	stringSlice := strings.Split(s, " ")
	cleanedStringSlice := []string{}
	for _, word := range stringSlice{ 
	
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax"{
			cleanedStringSlice = append(cleanedStringSlice, "****")
		}else{
			cleanedStringSlice = append(cleanedStringSlice, word)
		}
	}
	return strings.Join(cleanedStringSlice, " ")

}
