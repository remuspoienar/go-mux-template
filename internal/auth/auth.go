package auth

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"starter/internal/database"
	"strconv"
	"strings"
	"time"
)

func CreateToken(id int, jwtSecret string) (string, error) {
	jwtPayload := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		Subject:   strconv.Itoa(id),
	})

	token, err := jwtPayload.SignedString([]byte(jwtSecret))

	return token, err
}

func CreateRefreshToken() *database.RefreshToken {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	refreshToken := hex.EncodeToString(bytes)

	return &database.RefreshToken{
		Token: refreshToken,
		Exp:   time.Now().Add(time.Duration(60*24) * time.Hour),
	}

}

func ExtractTokenFromHeader(r *http.Request, jwtSecret string) (*jwt.RegisteredClaims, bool) {
	tokenStr := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")

	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if tokenStr == "" || err != nil {
		return nil, false
	}

	return token.Claims.(*jwt.RegisteredClaims), true
}
