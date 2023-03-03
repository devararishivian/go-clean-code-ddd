package service

type AuthService interface {
	Authenticate(email, password string) error
	Revoke(email string) error
}
