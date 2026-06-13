package middleware

import (
	"log"
	"net/http"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil dari konteks pake kunci "role" (HARUS SAMA PERSIS)
		val := r.Context().Value("role")

		log.Printf("DEBUG ADMIN - Role dari konteks: %v", val)

		role, ok := val.(string)
		if !ok || role != "admin" {
			log.Printf("DEBUG ADMIN - Akses ditolak! Role: %v", val)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
