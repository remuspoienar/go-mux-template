package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"starter/internal/database"
)

func (c *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	users := database.GetUsers()

	var id int
	if len(users) == 0 {
		id = 1
	} else if last := users[len(users)-1]; last.Id != 0 {
		id = last.Id + 1
	}
	body, _ := io.ReadAll(r.Body)
	p := userPayload{}
	u := database.User{}
	err := json.Unmarshal(body, &p)
	if err != nil {
		log.Fatal(err)
	}

	u.Id = id
	u.Email = p.Email
	u.SetPassword(p.Password)

	users = append(users, u)
	database.SaveUsers(users)

	respondWithJSON(w, 201, u.PublicFields())
}

func (c *ApiConfig) UpdateUser(w http.ResponseWriter, r *http.Request) {

	users := database.GetUsers()

	user := r.Context().Value("currentUser").(*database.User)

	body, _ := io.ReadAll(r.Body)
	p := userPayload{}
	err := json.Unmarshal(body, &p)
	if err != nil {
		log.Fatal(err)
	}

	if p.Email != "" {
		user.Email = p.Email
	}
	if p.Password != "" {
		user.SetPassword(p.Password)
	}

	database.SaveUsers(users)

	respondWithJSON(w, 200, user.PublicFields())
}
