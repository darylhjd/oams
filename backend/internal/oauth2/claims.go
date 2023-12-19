package oauth2

type Claims interface {
	UserID() string
	UserEmail() string
}
