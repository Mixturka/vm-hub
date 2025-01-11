package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Mixturka/vm-hub/internal/app/application/dtos"
	"github.com/Mixturka/vm-hub/internal/app/application/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var registerDto dtos.RegisterDto

	if err := json.NewDecoder(r.Body).Decode(&registerDto); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := ac.authService.ValidateDto(registerDto); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := ac.authService.Register(registerDto, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		// "user": map[string]interface{}{
		// 	"id":    user.ID,
		// 	"name":  user.Name,
		// 	"email": user.Email,
		// },
	})
}

func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginDto dtos.LoginDto

	if err := json.NewDecoder(r.Body).Decode(&loginDto); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := ac.authService.ValidateDto(loginDto); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := ac.authService.Login(loginDto, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successfull",
	})
}

func (ac *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	err := ac.authService.Logout(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Logout successfull",
	})
}
