package main

import (
	"blogBackend/internal/database"
	"blogBackend/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// load db conn
	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// initialize db
	database.InitDB(connStr)

	// set up router
	r := mux.NewRouter()

	r.HandleFunc("/blogs", handlers.GetBlogs).Methods("GET")
	r.HandleFunc("/blogs/{id}", handlers.GetBlog).Methods("GET")
	r.HandleFunc("/blogs", handlers.CreateBlog).Methods("POST")
	r.HandleFunc("/blogs/{id}", handlers.DeleteBlog).Methods("DELETE")

	// start server
	log.Println("Server starting on port: 8000...")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}