package entity

type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}
