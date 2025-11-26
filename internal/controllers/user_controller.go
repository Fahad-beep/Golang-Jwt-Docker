package controllers

import (
	"encoding/json"
	"log"
	"main/internal/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	Repo *models.UserRepository
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		http.Error(w, "Invalid Email or Password", http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Password couldnt be hashed", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashed)

	err = c.Repo.CreateUser(&user)
	if err != nil {
		http.Error(w, "Error saving user (user might already exists)", http.StatusInternalServerError)
		log.Println("Error while cerating user ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "user succesfully registered",
		"user_id": user.ID,
	})
}
