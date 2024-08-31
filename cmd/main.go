package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fanzru/bythen/internal/app/user/app"
	"github.com/fanzru/bythen/internal/app/user/model"
	"github.com/fanzru/bythen/internal/app/user/port"
	"github.com/fanzru/bythen/internal/app/user/port/genhttp"
	"github.com/fanzru/bythen/internal/app/user/repo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

// Middleware function to log request details
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Middleware function to set a request timeout
func timeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware to authenticate using JWT tokens
func authMiddleware(secretKey string) genhttp.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := authHeader[len("Bearer "):]
			claims := &model.JWTClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			type contextKey string

			const userIDKey contextKey = "userID"

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Function to load environment variables with default values
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load environment variables
	secretKey := getEnv("JWT_SECRET_KEY", "default-secret-key")
	dsn := getEnv("DATABASE_DSN", "user:password@tcp(localhost:3306)/dbname")

	// Initialize the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the UserRepository and UserService
	userRepo := repo.NewUserRepository(db)
	userService := app.NewUserService(userRepo, secretKey)
	userHandler := port.NewUserHandler(userService)

	// Create a new ServerInterface implementation
	serverInterface := &genhttp.ServerInterfaceWrapper{
		Handler: userHandler,
		HandlerMiddlewares: []genhttp.MiddlewareFunc{
			loggingMiddleware,
			timeoutMiddleware,
		},
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Wrap the server in an http.Handler
	handler := genhttp.HandlerFromMux(serverInterface, mux)

	// Serve the Swagger UI
	mux.Handle("/doc/swagger/", http.StripPrefix("/doc/swagger", http.FileServer(http.Dir("./docs/swagger"))))

	// Start the HTTP server
	addr := ":8080"
	log.Printf("Server is running on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("could not start server: %s", err)
	}
}
