package entity

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtClaim struct {
	BookID string `json:"book_id"`
	jwt.StandardClaims
}

func (claims *JwtClaim) NewToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, _ := token.SignedString([]byte("secret"))
	return t
}
func RefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
