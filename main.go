package main

import (
	"database/sql"
	"os"

	"github.com/RFreyAcona/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	db             *database.Queries
	fileserverHits atomic.Int32
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		platform = "prod"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		db:       dbQueries,
		platform: platform,
	}

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /api/healthz", healthzHandler)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	handler := http.FileServer(http.Dir(filepathRoot))
	serveMux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(handler)))
	server := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}
	log.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
