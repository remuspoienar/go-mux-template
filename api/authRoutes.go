package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"slices"
	"starter/internal/auth"
	"starter/internal/database"
	"strings"
	"time"
)

func (c *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	users := database.GetUsers()

	body, _ := io.ReadAll(r.Body)
	p := userPayload{}
	err := json.Unmarshal(body, &p)
	if err != nil {
		log.Fatal(err)
	}

	userIdx := slices.IndexFunc(users, func(user database.User) bool {
		return user.Email == p.Email
	})

	if userIdx == -1 {
		respondWithJSON(w, 401, "unauthorized")
		return
	}

	user := &users[userIdx]
	err = p.Authenticate(user)

	if err != nil {
		respondWithJSON(w, 401, "unauthorized")
		return
	}

	token, err := auth.CreateToken(user.Id, c.JwtSecret)

	if err != nil {
		respondWithJSON(w, 401, "unauthorized")
		return
	}

	user.RefreshToken = auth.CreateRefreshToken()
	database.SaveUsers(users)

	result := user.PublicFields()
	result["token"] = token

	respondWithJSON(w, 200, result)
}

func (c *ApiConfig) RefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenStr := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")

	users := database.GetUsers()

	userIdx := slices.IndexFunc(users, func(user database.User) bool {
		t := user.RefreshToken
		return t != nil && t.Token == tokenStr && t.Exp.After(time.Now())
	})

	if userIdx == -1 {
		respondWithJSON(w, 401, "unauthorized")
		return
	}

	user := &users[userIdx]
	token, err := auth.CreateToken(user.Id, c.JwtSecret)

	if err != nil {
		respondWithJSON(w, 401, "unauthorized")
		return
	}

	respondWithJSON(w, 200, map[string]any{"token": token})

}

func (c *ApiConfig) RevokeToken(w http.ResponseWriter, r *http.Request) {
	tokenStr := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")

	users := database.GetUsers()

	userIdx := slices.IndexFunc(users, func(user database.User) bool {
		t := user.RefreshToken
		return t.Token == tokenStr
	})

	if userIdx == -1 {
		respondWithJSON(w, 401, "unauthorized")
		return
	}

	user := &users[userIdx]

	user.RefreshToken = nil
	database.SaveUsers(users)
	respondWithJSON(w, 204, nil)
}
