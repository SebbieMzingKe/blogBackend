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
	// Extract user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch only blogs belonging to the logged-in user
	rows, err := database.DB.Query("SELECT id, title, body, author, created_at FROM blogs WHERE author = $1", userID)
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
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}


    err := database.DB.QueryRow(
        "INSERT INTO blogs (title, body, author) VALUES ($1, $2, $3) RETURNING id, created_at",
        b.Title, b.Body, userID,
    ).Scan(&b.ID, &b.CreatedAt)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// Extract authenticated user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if the blog belongs to the user
	var authorID int
	err = database.DB.QueryRow("SELECT author FROM blogs WHERE id = $1", id).Scan(&authorID)
	if err == sql.ErrNoRows {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error retrieving blog", http.StatusInternalServerError)
		return
	}

	// Ensure the authenticated user is the blog owner
	if authorID != userID {
		http.Error(w, "Forbidden: You can only delete your own blogs", http.StatusForbidden)
		return
	}

	// Proceed with deletion if the user owns the blog
	result, err := database.DB.Exec("DELETE FROM blogs WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error processing delete request", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
