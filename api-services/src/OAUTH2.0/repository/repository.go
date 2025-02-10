package repository

type Service interface {
	HandleGoogleLogin() (string, error)
	HandleGoogleCallback(code string) (*User, error)
	ValidateToken(token string) (*User, error)
}
