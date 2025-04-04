package handlers

import (
	"blogBackend/internal/database"
	"blogBackend/internal/middleware"
	"blogBackend/internal/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// get all blogs
func GetBlogs(w http.ResponseWriter, r *http.Request) {
    // Get email from context
    userEmail, ok := r.Context().Value("user_email").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Fetch blogs where the author is the logged-in user's email
    rows, err := database.DB.Query("SELECT id, title, body, author, created_at FROM blogs WHERE author = $1", userEmail)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var blogs []models.Blog
    for rows.Next() {
        var b models.Blog
        err := rows.Scan(&b.ID, &b.Title, &b.Body, &b.Author, &b.CreatedAt)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        blogs = append(blogs, b)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(blogs)
}

// get a single blog id
func GetBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Blog ID", http.StatusBadRequest)
		return
	}

	// Extract authenticated user ID
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var b models.Blog
	var authorID int

	// Fetch the blog and check ownership
	err = database.DB.QueryRow("SELECT id, title, body, author, created_at FROM blogs WHERE id = $1", id).
		Scan(&b.ID, &b.Title, &b.Body, &authorID, &b.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching blog", http.StatusInternalServerError)
		return
	}

	// Ensure the user is the owner of the blog
	if authorID != userID {
		http.Error(w, "Forbidden: You can only view your own blogs", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

// create new blog
func CreateBlog(w http.ResponseWriter, r *http.Request) {
    var b models.Blog
    if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Get email from context
    userEmail, ok := r.Context().Value("user_email").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Insert blog with associated email
    err := database.DB.QueryRow(
        "INSERT INTO blogs (title, body, author) VALUES ($1, $2, $3) RETURNING id, created_at",
        b.Title, b.Body, userEmail, 
    ).Scan(&b.ID, &b.CreatedAt)

    if err != nil {
        http.Error(w, "Error saving blog", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(b)
}


// delete a single blog
func DeleteBlog(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid Blog ID", http.StatusBadRequest)
        return
    }

    // Get logged-in user's email
    userEmail, ok := r.Context().Value("user_email").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Check if the blog belongs to the logged-in user
    var authorEmail string
    err = database.DB.QueryRow("SELECT author FROM blogs WHERE id = $1", id).Scan(&authorEmail)
    if err == sql.ErrNoRows {
        http.Error(w, "Blog not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Error retrieving blog", http.StatusInternalServerError)
        return
    }

    // Ensure the logged-in user is the owner
    if authorEmail != userEmail {
        http.Error(w, "Unauthorized: You can only delete your own blogs", http.StatusForbidden)
        return
    }

    // Delete the blog
    _, err = database.DB.Exec("DELETE FROM blogs WHERE id = $1", id)
    if err != nil {
        http.Error(w, "Error deleting blog", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

