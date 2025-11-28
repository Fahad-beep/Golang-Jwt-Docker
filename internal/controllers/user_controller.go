package controllers

import (
	"encoding/json"
	"log"
	"main/internal/auth"
	"main/internal/models"
	"net/http"
	"time"

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

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Print("credentials.Email: ", credentials.Email)
		log.Print("\ncredentials.Paassword: ", credentials.Password)
		log.Print("\nr.Body ", r.Body)
		http.Error(w, "Invalid Credentials (Re-enter)", http.StatusBadRequest)
		return
	}

	user, err := c.Repo.FetchUserByEmail(credentials.Email)
	if err != nil {
		http.Error(w, "Invalid Email or Password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, _ := auth.CreateAccessToken(user.ID, user.Email)
	refreshToken, _ := auth.GenerateRefereshToken()

	err = c.Repo.StoreRefreshTokens(user.ID, refreshToken, time.Now().Add(time.Hour*24*7))
	if err != nil {
		http.Error(w, "Server Error (failed to save hashed token in DB)", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "You are succesfully logged in",
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (c *UserController) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request (failed to send correct token to refresh)", http.StatusBadRequest)
		return
	}

	user, err := c.Repo.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid Or Session Expired", http.StatusUnauthorized)
		return
	}

	accessToken, _ := auth.CreateAccessToken(user.ID, user.Email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": accessToken,
	})
}
