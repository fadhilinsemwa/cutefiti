package main

import (
	"log"
	"net/http"
	"os"

	"usersb/handlers"
	"usersb/middleware"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the database handler
	dbHandler, err := handlers.NewDBHandler()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbHandler.CloseDB()

	// Initialize the ManufacturerHandler
	manufacturerHandler := handlers.NewManufacturerHandler(dbHandler.DB)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register manufacturer-related endpoints
	manufacturerHandler.RegisterManufacturerRoutes(mux)

	// Wrap the router with CORS middleware
	corsRouter := middleware.CORSMiddleware(mux)

	// Start the server
	port := os.Getenv("PORT2")
	if port == "" {
		port = "8002" // Default port if not specified
	}
	log.Printf("Server started on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsRouter))
}
