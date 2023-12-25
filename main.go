package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	jwtSecret      string
	polkaKey       string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable is not set")
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey,
	}

	router := chi.NewRouter()

	const filepathRoot = "."

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	api_router := chi.NewRouter()

	api_router.Get("/healthz", handlerReadiness)
	api_router.Post("/users", apiCfg.handlerUsersCreate)
	api_router.Post("/reservations", apiCfg.handlerReservationsCreate)

	api_router.Post("/login", apiCfg.handlerLogin)
	api_router.Post("/events", apiCfg.handlerEventsCreate)

	router.Mount("/api", api_router)

	corsMux := middlewareCors(router)

	srv := &http.Server{
		Handler: corsMux,
		Addr:    ":" + port,
	}

	fmt.Printf("Server starting on port %v\n", port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
