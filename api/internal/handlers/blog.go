package handlers

import (
	"blogBackend/internal/database"
	"blogBackend/internal/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// get all blogs
func GetBlogs(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query("SELECT id, title, body, author, created_at FROM blogs")
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
func GetBlog(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	var b models.Blog

	row := database.DB.QueryRow("SELECT id, title, body, author, created_at FROM blogs WHERE id = $1", id)
		row.Scan(&b.ID, &b.Title, &b.Body, &b.Author, &b.CreatedAt)

		if err == sql.ErrNoRows {
			http.Error(w, "Blog not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Error fetching blog", http.StatusInternalServerError)
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

    err := database.DB.QueryRow(
        "INSERT INTO blogs (title, body, author) VALUES ($1, $2, $3) RETURNING id, created_at",
        b.Title, b.Body, b.Author,
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
func DeleteBlog(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Invalid Blog ID", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("DELETE FROM blogs WHERE id = $1", id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		http.Error(w, "Error receiving delete command", http.StatusInternalServerError)
	}

	if rowsAffected == 0 {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}