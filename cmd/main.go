package main

import (
	"log"
	"net/http"
	"os"

	"github.com/amitp-vv/go-backend/internal/middleware"
	"github.com/amitp-vv/go-backend/internal/models"
	"github.com/amitp-vv/go-backend/internal/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func startEventListeners() {
	// TODO: Implement event listeners
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database
	if err := models.InitDb(); err != nil {
		log.Fatalf("DB init failed: %v", err)
	}

	// Start event listeners
	startEventListeners()

	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok": true}`))
	})

	// Apply middleware
	r.Use(middleware.RequireAuth)
	r.Use(middleware.RateLimit)

	// Set up routes
	routes.RegisterAdminRoutes(r)
	routes.RegisterAuthRoutes(r)
	routes.RegisterChainRoutes(r)
	routes.RegisterReadRoutes(r)

	// Start the server
	log.Printf("Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
