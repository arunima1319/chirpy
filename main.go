package main 

import (
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
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
	mux.HandleFunc("/healthz", readinessEndpoint)
	mux.HandleFunc("/metrics", apiCfg.displaysNumberOfRequests)
	mux.HandleFunc("/reset", apiCfg.resetHits)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func readinessEndpoint (w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
	
}

func(cfg *apiConfig) displaysNumberOfRequests (w http.ResponseWriter, r *http.Request){ 
	fmt.Fprintf(w, "Hits: %d", cfg.fileServerHits.Load())
}

func (cfg *apiConfig) resetHits (w http.ResponseWriter, r *http.Request){
	cfg.fileServerHits.Store(0)
	w.Write([]byte("has been reset"))
}