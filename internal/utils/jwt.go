package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(
	email string,
	name string,
	nim string,
) (string, error) {

	secretKey := os.Getenv("JWT_SECRET")

	if secretKey == "" {
		secretKey = "kunci-rahasia-pemira"
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"name":  name,
			"nim":   nim,
			"exp": time.Now().
				Add(24 * time.Hour).
				Unix(),
		},
	)

	return token.SignedString(
		[]byte(secretKey),
	)
}
