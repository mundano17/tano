package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"tanoserver/internal/auth"
	db "tanoserver/internal/db/generated"
	"tanoserver/internal/routes"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)
	queries.CreateUsersTable(ctx)
	authService := auth.NewAuthService(queries)
	authController := auth.NewAuthController(authService)

	mux := http.NewServeMux()
	mux.Handle("/user/", routes.NewAuthHandler(authController))
	mux.HandleFunc(
		"POST /refresh",
		authController.RefreshTokenController,
	)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Method:", r.Method)
		fmt.Println("Path:", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	})

	port := os.Getenv("SERVER_PORT")
	addr := fmt.Sprintf(":%s", port)

	fmt.Printf("Server is starting on port %s...\n", port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
