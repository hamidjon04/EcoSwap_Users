package model

type RefreshToken struct {
	UserId    string
	Token     string
	ExpiresAt int64
}

type InfoUser struct {
	Id       string
	Username string
	Password string
	FullName string
}

type ResetPassword struct {
	Email       string
	OldPassword string
	NewPassword string
}

type Error struct {
	Error string
}
