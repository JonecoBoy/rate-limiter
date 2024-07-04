package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JonecoBoy/rate-limiter/limiter"
	"github.com/JonecoBoy/rate-limiter/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	storageType := os.Getenv("RATE_LIMIT_STRATEGY")
	if storageType == "" {
		storageType = "memory"
	}

	rateLimiter := limiter.NewRateLimiter(storageType)

	r := mux.NewRouter()
	r.Use(rateLimiter.Middleware)

	r.HandleFunc("/", routes.HomeHandler).Methods("GET")

	webPort := os.Getenv("SERVER_PORT")

	server := &http.Server{
		Handler:      r,
		Addr:         ":" + webPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Server starting on port " + webPort)

	log.Fatal(server.ListenAndServe())
}
