package middleware

import (
	"log"
	"net/http"
)

// CORSMiddleware is a middleware function that handles CORS settings.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("CORS Middleware hit")

		// Log the incoming request method and headers
		log.Printf("Request Method: %s", r.Method)
		log.Printf("Request Headers: %v", r.Header)

		// Set CORS headers to allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight (OPTIONS) requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
