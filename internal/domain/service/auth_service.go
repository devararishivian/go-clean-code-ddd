package service

type AuthService interface {
	Authenticate(email, password string) error
	Revoke(email string) error
	GenerateToken(email string) (accessToken, refreshToken string, err error)
}
