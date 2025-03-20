package handlers

import (
	"blogBackend/internal/database"
	"blogBackend/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// get all blogs
func GetBlogs(w http.ResponseWriter, r *http.Request){
	rows, err := database.DB.Query("SELECT id, title, body, author, created_at FROM blogs")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	var blogs []models.Blog

	for rows.Next(){
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

	row := database.DB.QueryRow("SELECT id, title, body, author, created_at FROM blogs WHERE id = ?", id)
		row.Scan(&b.ID, &b.Title, &b.Body, &b.Author, &b.CreatedAt)

		if row != nil {
			http.Error(w, "Blog not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(b)
}

// create new blog
func CreateBlog(w http.ResponseWriter, r *http.Request)  {
	var b models.Blog

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
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

	result, err := database.DB.Exec("DELETE FROM blogs WHERE id = ?", id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}