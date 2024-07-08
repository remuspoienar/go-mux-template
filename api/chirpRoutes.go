package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"slices"
	"starter/internal/database"
	"strconv"
)

const (
	Asc  string = "asc"
	Desc        = "desc"
)

func CreateChirp(w http.ResponseWriter, r *http.Request) {
	chirps := database.GetChirps()

	user := r.Context().Value("currentUser").(*database.User)

	var id int
	if len(chirps) == 0 {
		id = 1
	} else if last := chirps[len(chirps)-1]; last.Id != 0 {
		id = last.Id + 1
	}
	body, _ := io.ReadAll(r.Body)
	c, err := validateChirp(body)
	if err != nil {
		respondWithJSON(w, 400, err)
		return
	}

	c.Id = id
	c.AuthorId = user.Id
	chirps = append(chirps, c)
	database.SaveChirps(chirps)

	respondWithJSON(w, 201, c)
}

func GetChirps(w http.ResponseWriter, r *http.Request) {
	chirps := database.GetChirps()
	authorIdParam := r.URL.Query().Get("author_id")
	if authorIdParam != "" {
		authorId, err := strconv.Atoi(authorIdParam)
		if err != nil {
			respondWithJSON(w, 400, "invalid author id")
			return
		}

		chirps = slices.DeleteFunc(chirps, func(chirp database.Chirp) bool {
			return chirp.AuthorId != authorId
		})
	}

	sortParam := r.URL.Query().Get("sort")
	if sortParam != "" {
		slices.SortFunc(chirps, func(a, b database.Chirp) int {
			if sortParam == Asc {
				return a.Id - b.Id
			}
			return b.Id - a.Id
		})
	}

	respondWithJSON(w, 200, chirps)
}

func GetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId, _ := strconv.Atoi(r.PathValue("id"))
	chirps := database.GetChirps()

	index := slices.IndexFunc(chirps, func(c database.Chirp) bool {
		return c.Id == chirpId
	})
	if index == -1 {
		respondWithJSON(w, 404, "Could not load chirp")
		return
	}

	respondWithJSON(w, 200, chirps[index])
}

func DeleteChirp(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("currentUser").(*database.User)

	chirpId, _ := strconv.Atoi(r.PathValue("id"))
	chirps := database.GetChirps()

	index := slices.IndexFunc(chirps, func(c database.Chirp) bool {
		return c.Id == chirpId
	})
	if index == -1 {
		respondWithJSON(w, 404, "Could not load chirp")
		return
	}

	chirp := &chirps[index]

	if chirp.AuthorId != user.Id {
		respondWithJSON(w, 403, "Only the author can delete this tweet")
		return
	}

	slices.DeleteFunc(chirps, func(chirp database.Chirp) bool {
		return chirp.Id == chirpId
	})
	database.SaveChirps(chirps)

	respondWithJSON(w, 204, nil)
}

func validateChirp(body []byte) (database.Chirp, error) {
	c := database.Chirp{}
	err := json.Unmarshal(body, &c)
	if err != nil {
		log.Fatal(err)
	}

	if len(c.Body) > 140 {
		return c, errors.New("max length is 140")
	}

	containsBadWords := regexp.MustCompile(`(?i) (kerfuffle|sharbert|fornax) `)
	c.Body = containsBadWords.ReplaceAllString(c.Body, " **** ")

	return c, nil

}
