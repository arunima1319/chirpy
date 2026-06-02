package main 

import (
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
	"encoding/json"
	"strings"
)

type apiConfig struct{ 
	fileServerHits atomic.Int32
}
func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", 
			apiCfg.middleWareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", readinessEndpoint)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.displaysNumberOfRequests)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHits)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request){ 

	type chirp struct{ 
		Body string `json:"body"`
	}

	type responseError struct{ 
		Error string `json:"error"`
	}
	type responseValid struct{
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirped := chirp{}
	err := decoder.Decode(&chirped)
	if err!= nil{ 
		log.Printf("Error in decoding JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(chirped.Body) > 140{ 
		respBody := responseError{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err!= nil{ 
			log.Printf("Error in Marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	cleanedChirp := replaceBadWord(chirped.Body)

	respBody := responseValid{ 
		CleanedBody: cleanedChirp,
	}
	dat, err := json.Marshal(respBody)
	if err!= nil{
		log.Printf("Error in marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	
	w.Header().Set("Content=Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

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
func readinessEndpoint (w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
	
}

func(cfg *apiConfig) displaysNumberOfRequests (w http.ResponseWriter, r *http.Request){ 
	w.Header().Set("Content-Type", "text/html")
	stringToReturn := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits.Load())
	w.Write([]byte(stringToReturn))
	//fmt.Fprintf(w, "Hits: %d", cfg.fileServerHits.Load())
}

func (cfg *apiConfig) resetHits (w http.ResponseWriter, r *http.Request){
	cfg.fileServerHits.Store(0)
	w.Write([]byte("has been reset"))
}