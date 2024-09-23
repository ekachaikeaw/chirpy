package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ekachaikeaw/chirpy/database"
)

var db *database.DB

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>
k
</html>`, cfg.fileserverHits)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hit counter reset"))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

// Helper function to respond with error
func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// Helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// Function to clean the body from profane words
func cleanProfanity(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		for _, profane := range profaneWords {
			if lowerWord == profane {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	// Decode the request body
	var params parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil || params.Body == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request: 'body' field is required")
		return
	}

	// Validate Chirp length
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Clean the body from profane words
	cleanedBody := cleanProfanity(params.Body)
	chirp, err := db.CreateChirp(cleanedBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
		return
	}
	
	// Respond with cleaned body
	respondWithJSON(w, http.StatusOK, chirp)
}

func getChirpys(w http.ResponseWriter, r *http.Request) {
    chirps, err := db.GetChirpy()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func main() {
	apiCfg := apiConfig{}
	var err error
	db, err = database.NewDB("database.json")
	if err != nil {
		fmt.Println("Error initializing database:", db)
	}
	fmt.Println("Connected to database")

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080", // กำหนดพอร์ต
		Handler: mux,     // ใช้ Handler แบบกำหนดเอง
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/chirps", validateChirp)
	mux.HandleFunc("GET /api/chirps", getChirpys)

	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
