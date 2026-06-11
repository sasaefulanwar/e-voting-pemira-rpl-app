package middleware

import (
	"context"
	"encoding/json"
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

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			println("🚨 SENSOR SATPAM GAGAL: ", err.Error())

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Token JWT lu gak valid atau kadaluarsa!"})
			return
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

		email := claims["email"].(string)

		ctx := context.WithValue(r.Context(), "email", email)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
