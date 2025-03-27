package main

import (
	"blogBackend/internal/database"
	"blogBackend/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load .env file (for local development)
	_ = godotenv.Load()

	// Get database connection string
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Initialize database
	if err := database.InitDB(connStr); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Set up router
	r := mux.NewRouter()

	// Default homepage
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🚀 Welcome to the Blog API! Access the API at /blogs")
	})

	// API routes
	r.HandleFunc("/blogs", handlers.GetBlogs).Methods("GET", "OPTIONS")
	r.HandleFunc("/blogs/{id}", handlers.GetBlog).Methods("GET", "OPTIONS")
	r.HandleFunc("/blogs", handlers.CreateBlog).Methods("POST", "OPTIONS")
	r.HandleFunc("/blogs/{id}", handlers.DeleteBlog).Methods("DELETE", "OPTIONS")

	// Enable CORS (Correct Configuration)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://react-blog-three-orcin.vercel.app", "http://localhost:3000"}, // React Frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap the router with CORS middleware
	handler := corsHandler.Handler(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port
	}
	log.Printf("Server running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
