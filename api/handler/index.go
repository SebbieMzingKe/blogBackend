package handler

import (
	"blogBackend/internal/database"
	"blogBackend/internal/handlers"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	_ = godotenv.Load()

	// load db conn
	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		http.Error(w, "DATABASE_URL environment variable not set", http.StatusInternalServerError)
	}

	if database.DB == nil {
		err := database.InitDB(connStr)
		if err != nil {
			http.Error(w, "Failed to initialize database", http.StatusInternalServerError)
			return
		}
	}

	// Set up router
	rtr := mux.NewRouter()

	// Default homepage
	rtr.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🚀 Welcome to the Blog API! Access the API at /api/blogs")
	})

	// API routes
	rtr.HandleFunc("/api/blogs", handlers.GetBlogs).Methods("GET")
	rtr.HandleFunc("/api/blogs/{id}", handlers.GetBlog).Methods("GET")
	rtr.HandleFunc("/api/blogs", handlers.CreateBlog).Methods("POST")
	rtr.HandleFunc("/api/blogs/{id}", handlers.DeleteBlog).Methods("DELETE")

	// Serve the request using the router
	rtr.ServeHTTP(w, r)
}
