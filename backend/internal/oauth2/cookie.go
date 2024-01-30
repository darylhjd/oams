package oauth2

import (
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	SessionCookieIdent = "oams_session_cookie"
)

func SetCookie(w http.ResponseWriter, token *oauth2.Token) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieIdent,
		Value:    token.AccessToken,
		Path:     "/",
		Expires:  token.Expiry,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
}

func DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieIdent,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
}
