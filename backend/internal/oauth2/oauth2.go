package oauth2

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/darylhjd/oams/backend/internal/env"
)

const (
	SessionCookieIdent = "oams_session_cookie"
)

var (
	MicrosoftAuthority = fmt.Sprintf("https://login.microsoftonline.com/%s/", env.GetAPIServerAzureTenantID())
	KeySetSource       = fmt.Sprintf("https://login.microsoftonline.com/%s/discovery/v2.0/keys", env.GetAPIServerAzureTenantID())
)

var (
	tokenIssuer = MicrosoftAuthority + "v2.0"
)

// AzureClaims is a custom struct to hold the claims received from Microsoft Azure AD.
type AzureClaims struct {
	jwt.RegisteredClaims
	Scp   string
	Roles []string
}

// CheckAzureToken to make sure the token passes validation.
func CheckAzureToken(set jwk.Set, tokenString string) (*AzureClaims, *jwt.Token, error) {
	// https://learn.microsoft.com/en-us/azure/active-directory/develop/access-tokens#validating-tokens
	claims := &AzureClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Header["typ"] != "JWT" {
			return nil, errors.New("wrong token type")
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid not found")
		}

		keys, ok := set.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %s not found", kid)
		}

		publicKey := &rsa.PublicKey{}
		if err := keys.Raw(publicKey); err != nil {
			return nil, errors.New("could not parse key")
		}

		return publicKey, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}),
		jwt.WithAudience(env.GetAPIServerAzureClientID()),
		jwt.WithIssuer(tokenIssuer))
	if err != nil || !token.Valid {
		return nil, nil, fmt.Errorf("token check failed: %w", err)
	}

	return claims, token, nil
}

// SetSessionCookie to a response. Returns a copy of the session cookie that was set.
func SetSessionCookie(w http.ResponseWriter, res confidential.AuthResult) http.Cookie {
	// We use the account's home account ID as the key value identifier.
	cookie := &http.Cookie{
		Name:     SessionCookieIdent,
		Value:    res.Account.HomeAccountID,
		Path:     "/api/v1/",
		Expires:  res.ExpiresOn,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)
	return *cookie
}

// DeleteSessionCookie sets an expired OAMS session cookie to a response to request a browser to delete
// the current session cookie.
func DeleteSessionCookie(w http.ResponseWriter) http.Cookie {
	cookie := &http.Cookie{
		Name:     SessionCookieIdent,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)
	return *cookie
}
