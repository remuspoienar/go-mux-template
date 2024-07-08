package api

import (
	"context"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"slices"
	"starter/internal/auth"
	"starter/internal/database"
	"strconv"
)

type ApiConfig struct {
	FileserverHits int
	JwtSecret      string
	PolkaApiKey    string
}

type userPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Exp      int    `json:"expires_in_seconds"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func (p *userPayload) Authenticate(u *database.User) error {
	return bcrypt.CompareHashAndPassword(u.Hash, []byte(p.Password))
}

func (c *ApiConfig) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users := database.GetUsers()

		var user *database.User
		var userIdx int
		if claims, ok := auth.ExtractTokenFromHeader(r, c.JwtSecret); ok {
			userIdx = slices.IndexFunc(users, func(user database.User) bool {
				uId, _ := claims.GetSubject()
				userId, _ := strconv.Atoi(uId)
				return user.Id == userId
			})

			if userIdx == -1 {
				respondWithJSON(w, 401, "unauthorized")
				return
			}

			user = &users[userIdx]
		} else {
			respondWithJSON(w, 401, "unauthorized")
			return
		}
		ctx := context.WithValue(r.Context(), "currentUser", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
