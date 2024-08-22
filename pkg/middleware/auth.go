package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

type AuthHandler struct {
	Username string
	Password string
}

func (ah *AuthHandler) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(ah.Username))
			expectedPasswordHash := sha256.Sum256([]byte(ah.Password))

			if subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1 &&
				subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1 {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
