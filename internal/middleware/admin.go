package middleware

import (
	"net/http"
)

func AdminOnly(next http.Handler) http.Handler {

	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			roleVal :=
				r.Context().Value("role")

			if roleVal == nil {

				http.Error(
					w,
					"unauthorized",
					http.StatusUnauthorized,
				)

				return
			}

			role, ok :=
				roleVal.(string)

			if !ok || role != "admin" {

				http.Error(
					w,
					"forbidden: admin only",
					http.StatusForbidden,
				)

				return
			}

			next.ServeHTTP(w, r)
		},
	)
}
