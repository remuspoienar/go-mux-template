package api

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"starter/internal/database"
	"strings"
)

const eventUserUpgraded = "user.upgraded"

type whBody struct {
	Event string `json:"event"`
	Data  struct {
		UserId int `json:"user_id"`
	} `json:"data"`
}

func (c *ApiConfig) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := strings.ReplaceAll(strings.TrimSpace(r.Header.Get("Authorization")), "ApiKey ", "")

	if apiKey != c.PolkaApiKey {
		respondWithJSON(w, 401, "unauthorized")
		return
	}
	body, _ := io.ReadAll(r.Body)

	var wh whBody
	json.Unmarshal(body, &wh)

	if wh.Event != eventUserUpgraded {
		respondWithJSON(w, 204, nil)
		return
	}

	userId := wh.Data.UserId
	users := database.GetUsers()
	index := slices.IndexFunc(users, func(user database.User) bool {
		return user.Id == userId
	})

	if index == -1 {
		respondWithJSON(w, 404, "user not found")
		return
	}

	user := &users[index]
	user.IsChirpyRed = true
	database.SaveUsers(users)

	respondWithJSON(w, 204, nil)

}
