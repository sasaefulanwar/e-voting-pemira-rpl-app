package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type AuthService interface {
	GetGoogleLoginURL() string
	ProcessGoogleCallback(code string) (string, error)
}

type authService struct {
	oauthConf *oauth2.Config
}

func NewAuthService(oauthConf *oauth2.Config) AuthService {
	return &authService{oauthConf: oauthConf}
}

func (s *authService) GetGoogleLoginURL() string {
	return s.oauthConf.AuthCodeURL("state-pemira")
}

func (s *authService) ProcessGoogleCallback(code string) (string, error) {
	ctx := context.Background()

	token, err := s.oauthConf.Exchange(ctx, code)
	if err != nil {
		return "", errors.New("gagal menukar kode dari Google")
	}

	client := s.oauthConf.Client(ctx, token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || res.StatusCode != http.StatusOK {
		return "", errors.New("gagal mengambil profil Google")
	}
	defer res.Body.Close()

	var googleUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(res.Body).Decode(&googleUser); err != nil {
		return "", errors.New("gagal membaca profil Google")
	}

	admins := strings.Split(
		os.Getenv("ADMIN_EMAILS"),
		",",
	)

	role := "voter"

	for _, adminEmail := range admins {

		if strings.EqualFold(
			strings.TrimSpace(adminEmail),
			googleUser.Email,
		) {

			role = "admin"
			break
		}
	}

	jwtToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": googleUser.Email,
			"name":  googleUser.Name,
			"nim":   "",
			"role":  role,
			"exp":   time.Now().Add(24 * time.Hour).Unix(),
		},
	)

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "kunci-rahasia-pemira"
	}

	tokenString, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.New("gagal membuat token sesi")
	}

	return tokenString, nil
}
