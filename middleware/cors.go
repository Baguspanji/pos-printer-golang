package middleware

import (
	"log"
	"net/http"
	"time"
)

// CORS adds CORS headers to all responses and logs incoming requests
func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log incoming request
		log.Printf("[%s] %s %s - Remote: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			r.RemoteAddr)

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next(w, r)
	}
}
