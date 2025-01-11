package middleware

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Mixturka/vm-hub/internal/app/config"
)

type RecaptchaResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes,omitempty"`
}

func RecaptchaMiddleware(cfg *config.GRecapOptions, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recaptchaToken := r.Header.Get("recaptcha")
		if recaptchaToken == "" {
			http.Error(w, "Missing reCAPTCHA token", http.StatusBadRequest)
			return
		}

		secretKey := cfg.SecretKey
		if secretKey == "" {
			http.Error(w, "Server misconfiguration: Missing reCAPTCHA secret key", http.StatusInternalServerError)
			return
		}

		recaptchaURL := "https://www.google.com/recaptcha/api/siteverify"
		response, err := http.PostForm(recaptchaURL, map[string][]string{
			"secret":   {secretKey},
			"response": {recaptchaToken},
		})
		if err != nil {
			http.Error(w, "Failed to verify reCAPTCHA", http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			http.Error(w, "Failed to read reCAPTCHA verification response", http.StatusInternalServerError)
			return
		}

		var recaptchaResponse RecaptchaResponse
		if err := json.Unmarshal(body, &recaptchaResponse); err != nil {
			http.Error(w, "Failed to parse reCAPTCHA verification response", http.StatusInternalServerError)
			return
		}

		if !recaptchaResponse.Success {
			http.Error(w, "Invalid  reCAPTCHA token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
