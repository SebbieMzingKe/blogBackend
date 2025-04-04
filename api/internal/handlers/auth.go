package handlers

import (
	"blogBackend/internal/database"
	"blogBackend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secret-key")

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Credentials did not match", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	if err := database.CreateUser(user.Email, string(hashedPassword), user.Name); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user created"})
}

func SignIn(w http.ResponseWriter, r *http.Request)  {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error logging in", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByEmail(req.Email)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		fmt.Println("error logging user", err)
		return
	}


	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email": user.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey); 
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return 
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	// Remove the token from cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Expire immediately
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
