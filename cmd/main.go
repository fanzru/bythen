package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	blogApp "github.com/fanzru/bythen/internal/app/blog/app"
	blogPort "github.com/fanzru/bythen/internal/app/blog/port"
	blogGenhttp "github.com/fanzru/bythen/internal/app/blog/port/genhttp"
	blogRepo "github.com/fanzru/bythen/internal/app/blog/repo"
	userApp "github.com/fanzru/bythen/internal/app/user/app"
	"github.com/fanzru/bythen/internal/app/user/model"
	userPort "github.com/fanzru/bythen/internal/app/user/port"
	userGenhttp "github.com/fanzru/bythen/internal/app/user/port/genhttp"
	userRepo "github.com/fanzru/bythen/internal/app/user/repo"

	commentApp "github.com/fanzru/bythen/internal/app/comment/app"
	commentPort "github.com/fanzru/bythen/internal/app/comment/port"
	commentGenhttp "github.com/fanzru/bythen/internal/app/comment/port/genhttp"
	commentRepo "github.com/fanzru/bythen/internal/app/comment/repo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

// Middleware function to log request details
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r.WithContext(r.Context()))
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
func authMiddleware(secretKey string) func(http.Handler) http.Handler {
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

			const userIDKey string = "userID"
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

// Middleware to handle CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If the request is for OPTIONS, stop there
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
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
	userRepo := userRepo.NewUserRepository(db)
	userService := userApp.NewUserService(userRepo, secretKey)
	userHandler := userPort.NewUserHandler(userService)

	// Initialize the BlogRepository and BlogService
	blogRepo := blogRepo.NewPostRepository(db)
	blogService := blogApp.NewPostService(blogRepo)
	blogHandler := blogPort.NewPostHandler(blogService)

	// Initialize the CommentRepository and CommentService
	commentRepo := commentRepo.NewCommentRepository(db)
	commentService := commentApp.NewCommentService(commentRepo)
	commentHandler := commentPort.NewCommentHandler(commentService)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register the user routes with middlewares applied
	userGenhttp.HandlerWithOptions(userHandler, userGenhttp.StdHTTPServerOptions{
		BaseRouter: mux,
		Middlewares: []userGenhttp.MiddlewareFunc{
			loggingMiddleware,
			timeoutMiddleware,
		},
	})

	// Register the blog routes with middlewares including auth middleware applied
	blogGenhttp.HandlerWithOptions(blogHandler, blogGenhttp.StdHTTPServerOptions{
		BaseRouter: mux,
		Middlewares: []blogGenhttp.MiddlewareFunc{
			loggingMiddleware,
			timeoutMiddleware,
			authMiddleware(secretKey),
		},
	})

	// Register the comment routes with middlewares including auth middleware applied
	commentGenhttp.HandlerWithOptions(commentHandler, commentGenhttp.StdHTTPServerOptions{
		BaseRouter: mux,
		Middlewares: []commentGenhttp.MiddlewareFunc{
			loggingMiddleware,
			timeoutMiddleware,
			authMiddleware(secretKey),
		},
	})

	// Apply CORS middleware
	corsMux := corsMiddleware(mux)

	// Serve the Swagger UI
	mux.Handle("/doc/swagger/", http.StripPrefix("/doc/swagger", http.FileServer(http.Dir("/app/docs/swagger"))))

	// Start the HTTP server
	addr := ":8080"
	log.Printf("Server is running on %s", addr)
	if err := http.ListenAndServe(addr, corsMux); err != nil {
		log.Fatalf("could not start server: %s", err)
	}
}
