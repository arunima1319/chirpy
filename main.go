package main 

import (
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"os"
	"database/sql"
	"github.com/arunima1319/chirpy/internal/database"
	
)

type apiConfig struct{ 
	fileServerHits atomic.Int32
	dbQueries *database.Queries
	platform string
	secret string
}
func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	platformCurrent := os.Getenv("PLATFORM")
	secretForJWT := os.Getenv("SECRET")


	db, err := sql.Open("postgres", dbUrl)
	if err!= nil{
		log.Printf("Error in opening the database: %s", err)
		return
	}
	
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{}
	apiCfg.dbQueries = database.New(db)
	apiCfg.platform = platformCurrent
	apiCfg.secret = secretForJWT




	mux := http.NewServeMux()
	mux.Handle("/app/", 
			apiCfg.middleWareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", readinessEndpoint)
	//HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateDetails)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetOneChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	mux.HandleFunc("GET /admin/metrics", apiCfg.displaysNumberOfRequests)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetDatabase)

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
	w.Header().Set("Content-Type", "text/html")
	stringToReturn := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits.Load())
	w.Write([]byte(stringToReturn))
	//fmt.Fprintf(w, "Hits: %d", cfg.fileServerHits.Load())
}

func (cfg *apiConfig) resetHits (w http.ResponseWriter, r *http.Request){
	cfg.fileServerHits.Store(0)
	w.Write([]byte("has been reset"))
}