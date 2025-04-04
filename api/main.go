package main

import (
	"blogBackend/internal/database"
	"blogBackend/internal/handlers"
	"blogBackend/internal/middleware"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

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

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🚀 Welcome to the Blog API! Access the API at /blogs")
	})

	// API routes
	r.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	r.HandleFunc("/signin", handlers.SignIn).Methods("POST")
	r.Handle("/logout", middleware.AuthMiddleware(http.HandlerFunc(handlers.SignOut))).Methods("POST")
	r.Handle("/blogs", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetBlogs))).Methods("GET", "OPTIONS")
	r.Handle("/blogs/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetBlog))).Methods("GET", "OPTIONS")
	r.Handle("/blogs", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateBlog))).Methods("POST", "OPTIONS")
	r.Handle("/blogs/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteBlog))).Methods("DELETE", "OPTIONS")
	

	// Enable CORS 
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://react-blog-three-orcin.vercel.app", "http://localhost:3000"},
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
