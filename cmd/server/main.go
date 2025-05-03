package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fallenkarma/wasatext/internal/handlers"
	"github.com/fallenkarma/wasatext/internal/repository/postgres"
	"github.com/fallenkarma/wasatext/internal/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize repository
	// This is just a placeholder - you'll need to replace this with your actual repository implementation

	err := godotenv.Load()
    if err != nil {
        log.Println("Warning: Error loading .env file:", err)
    }

    // Get environment variables
    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080" // Default port if not specified
    }
    
    dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	log.Println("Warning: Error loading .env file:", err)

	repo,err := postgres.NewPostgresRepository(dbConnectionString,"/internal/images")
	if err != nil {
		log.Fatalf("Connection to database failed: %v", err)
	}

	// Initialize service with repository
	svc := service.New(repo)

	// Initialize handlers with service
	handler := handlers.New(svc)

	// Initialize router
	r := mux.NewRouter()
	
	// Add API prefix
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Public routes (no auth required)
	apiRouter.HandleFunc("/session", handler.Login).Methods("POST")

	// Protected routes (auth required)
	protected := apiRouter.NewRoute().Subrouter()
	protected.Use(handler.AuthMiddleware)

	// User routes
	protected.HandleFunc("/users/me/username", handler.SetMyUserName).Methods("PUT")
	protected.HandleFunc("/users/me/photo", handler.SetMyPhoto).Methods("PUT")

	// Conversation routes
	protected.HandleFunc("/conversations", handler.GetMyConversations).Methods("GET")
	protected.HandleFunc("/conversations/{id}", handler.GetConversation).Methods("GET")

	// Message routes
	protected.HandleFunc("/messages", handler.SendMessage).Methods("POST")
	protected.HandleFunc("/messages/forward", handler.ForwardMessage).Methods("POST")
	protected.HandleFunc("/messages/{id}/comment", handler.CommentMessage).Methods("POST")
	protected.HandleFunc("/messages/{id}/comment", handler.UncommentMessage).Methods("DELETE")
	protected.HandleFunc("/messages/{id}", handler.DeleteMessage).Methods("DELETE")

	// Group routes
	protected.HandleFunc("/groups/{id}/members", handler.AddToGroup).Methods("POST")
	protected.HandleFunc("/groups/{id}/leave", handler.LeaveGroup).Methods("POST")
	protected.HandleFunc("/groups/{id}/name", handler.SetGroupName).Methods("PUT")
	protected.HandleFunc("/groups/{id}/photo", handler.SetGroupPhoto).Methods("PUT")

	// Create server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on :" + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}