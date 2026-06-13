package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt_token")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Lu belum login Google cuy!"})
			return
		}

		tokenString := cookie.Value
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			secretKey = "kunci-rahasia-pemira"
		}

		log.Println("COOKIE FOUND:", true)

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			},
		)

		log.Println("JWT PARSE ERROR:", err)

		if token != nil {
			log.Println("TOKEN VALID:", token.Valid)
		}

		if err != nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Token JWT lu gak valid atau kadaluarsa!"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Gagal baca data token"})
			return
		}

		email, _ := claims["email"].(string)
		nim, _ := claims["nim"].(string)
		role, _ := claims["role"].(string)

		ctx := context.WithValue(
			r.Context(),
			"email",
			email,
		)

		ctx = context.WithValue(
			ctx,
			"nim",
			nim,
		)

		ctx = context.WithValue(
			ctx,
			"role",
			role,
		)

		log.Println("PATH:", r.URL.Path)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
