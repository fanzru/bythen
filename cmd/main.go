package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fanzru/bythen/internal/app/user/app"
	"github.com/fanzru/bythen/internal/app/user/model"
	"github.com/fanzru/bythen/internal/app/user/port"
	"github.com/fanzru/bythen/internal/app/user/port/genhttp"
	"github.com/fanzru/bythen/internal/app/user/repo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// Middleware function to log request details
func loggingMiddleware(next nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)
		return next(ctx, w, r, request)
	}
}

// Middleware function to set a request timeout
func timeoutMiddleware(next nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		return next(ctx, w, r.WithContext(ctx), request)
	}
}

// Middleware to authenticate using JWT tokens
func authMiddleware(secretKey string) nethttp.StrictHTTPMiddlewareFunc {
	return func(next nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return nil, nil
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
				return nil, nil
			}

			ctx = context.WithValue(ctx, "userID", claims.UserID)

			return next(ctx, w, r, request)
		}
	}
}

func main() {
	// Load the JWT secret key from the environment
	// secretKey := os.Getenv("JWT_SECRET_KEY")
	// if secretKey == "" {
	// 	log.Fatal("JWT_SECRET_KEY environment variable is not set")
	// }

	secretKey := "secretkey"

	// Initialize the database connection
	// dsn := "your-dsn-here" // replace with your actual DSN
	// db, err := sql.Open("mysql", dsn)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// Initialize the UserRepository and UserService
	userRepo := repo.NewUserRepository(nil)
	userService := app.NewUserService(userRepo, secretKey)
	userHandler := port.NewUserHandler(userService)

	// Create a new StrictHandler with middlewares
	strictServer := genhttp.NewStrictHandler(userHandler, []nethttp.StrictHTTPMiddlewareFunc{
		loggingMiddleware,
		timeoutMiddleware,
		authMiddleware(secretKey), // Apply authentication middleware
	})

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Wrap the strict server in an http.Handler
	handler := genhttp.HandlerWithOptions(strictServer, genhttp.StdHTTPServerOptions{})

	// Serve the Swagger UI
	mux.Handle("/doc/swagger/", http.StripPrefix("/doc/swagger", http.FileServer(http.Dir("./docs/swagger"))))

	// Add the strict server to the mux
	mux.Handle("/api/", handler)

	// Start the HTTP server
	addr := ":8080"
	log.Printf("Server is running on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("could not start server: %s", err)
	}
}
