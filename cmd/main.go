package main

import (
	"log"
	"net/http"
)

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Serve Swagger UI assets
	swaggerFS := http.FileServer(http.Dir("./docs/swagger/asset"))
	mux.Handle("/doc/swagger/asset/", http.StripPrefix("/doc/swagger/asset", swaggerFS))

	// Serve the Swagger JSON directly
	docFS := http.FileServer(http.Dir("./docs/swagger"))
	mux.Handle("/doc/swagger.json", ChainHandler(
		http.StripPrefix("/doc", docFS),
		SetResponseHeader("Access-Control-Allow-Origin", "*"),
	))

	// Serve the Swagger UI (index.html)
	mux.Handle("/doc/swagger/", http.StripPrefix("/doc/swagger", docFS))

	// Start the HTTP server
	addr := ":8080"
	log.Printf("Server is running on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("could not start server: %s", err)
	}
}

// ChainHandler chains middleware with the final handler
func ChainHandler(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

// SetResponseHeader sets custom headers for the HTTP response
func SetResponseHeader(key, value string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, value)
			next.ServeHTTP(w, r)
		})
	}
}
