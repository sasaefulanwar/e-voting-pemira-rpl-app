package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"pemira-rpl/internal/service"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	srv service.AuthService
}

func NewAuthHandler(srv service.AuthService) *AuthHandler {
	return &AuthHandler{srv: srv}
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.srv.GetGoogleLoginURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Kode OAuth kosong"})
		return
	}

	jwtToken, err := h.srv.ProcessGoogleCallback(code)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "jwt_token",
			Value:    jwtToken,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			MaxAge:   86400,
			Expires:  time.Now().Add(24 * time.Hour),
		},
	)

	http.Redirect(
		w,
		r,
		"http://localhost:5173/bind",
		http.StatusTemporaryRedirect,
	)
}

func (h *AuthHandler) BindNIM(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt_token")
	if err != nil {
		http.Error(w, "JWT tidak ditemukan", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "JWT tidak valid", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "JWT claim error", http.StatusUnauthorized)
		return
	}

	var body struct {
		NIM string `json:"nim"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	claims["nim"] = body.NIM

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, _ := newToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    jwtToken,
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"bind berhasil, JWT diperbarui"}`))
}

func (h *AuthHandler) Logout(
	w http.ResponseWriter,
	r *http.Request,
) {

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "jwt_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		},
	)

	

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "logout berhasil",
		},
	)
}

func (h *AuthHandler) Me(
	w http.ResponseWriter,
	r *http.Request,
) {

	emailValue :=
		r.Context().
			Value("email")

	roleValue :=
		r.Context().
			Value("role")

	nimValue :=
		r.Context().
			Value("nim")

	if emailValue == nil {

		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)

		return
	}

	json.NewEncoder(
		w,
	).Encode(
		map[string]interface{}{
			"email": emailValue,
			"role":  roleValue,
			"nim":   nimValue,
		},
	)
}
