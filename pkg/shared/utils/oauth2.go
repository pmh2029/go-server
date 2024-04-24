package utils

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const OauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func SetupConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "507025329306-07bihn57phj6t1c750ahea66t5v0mjsm.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/api/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GenerateStateOauthCookie(c *gin.Context) string {
	var expiration = time.Now().Add(2 * time.Minute)
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)

	return state
}
