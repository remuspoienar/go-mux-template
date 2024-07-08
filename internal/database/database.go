package database

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"sync"
	"time"
)

type Chirp struct {
	Body     string `json:"body"`
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
}

type RefreshToken struct {
	Token string    `json:"token"`
	Exp   time.Time `json:"exp"`
}

type User struct {
	Id           int           `json:"id"`
	Email        string        `json:"email"`
	Hash         []byte        `json:"hash,string"`
	IsChirpyRed  bool          `json:"is_chirpy_red"`
	RefreshToken *RefreshToken `json:"refresh_token,omitempty"`
}

type DB struct {
	Chirps []Chirp `json:"Chirps"`
	Users  []User  `json:"Users"`
}

func NewDB() *DB {
	return &DB{Chirps: []Chirp{}, Users: []User{}}
}

func ResetDb() {
	var dbg bool
	flag.BoolVar(&dbg, "debug", false, "Enable debug mode")
	flag.Parse()
	if dbg {
		err := os.Remove(path)

		if err != nil {
			fmt.Println("database.json already removed, skipping...")
		} else {
			fmt.Println("Deleted database.json")
		}
	}
}

func (u *User) SetPassword(pass string) {
	u.Hash, _ = bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
}

func (u *User) PublicFields() map[string]any {
	base := map[string]any{
		"id":            u.Id,
		"email":         u.Email,
		"is_chirpy_red": u.IsChirpyRed,
	}

	if u.RefreshToken != nil {
		base["refresh_token"] = u.RefreshToken.Token
	}

	return base
}

func GetUsers() []User {
	return unmarshal().Users
}

func SaveUsers(users []User) {
	old := unmarshal()
	old.Users = users
	marshal(old)
}

func GetChirps() []Chirp {
	return unmarshal().Chirps
}

func SaveChirps(chirps []Chirp) {
	old := unmarshal()
	old.Chirps = chirps
	marshal(old)
}

const path = "database.json"

var mu = &sync.RWMutex{}

func unmarshal() *DB {
	db := NewDB()
	_, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)

	contents, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if len(contents) == 0 {
		contents = []byte(`{"users": [], "chirps": []}`)
	}

	json.Unmarshal(contents, db)

	return db
}

func marshal(data *DB) {
	mu.Lock()
	defer mu.Unlock()
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	json, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(json)
	if err != nil {
		log.Fatal(err)
	}

}
