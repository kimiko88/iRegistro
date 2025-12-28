package domain

import "context"

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	GetByExternalID(ctx context.Context, externalID string) (*User, error)
	FindAll(schoolID uint) ([]User, error)
	Update(user *User) error
}

type AuthRepository interface {
	CreateSession(session *Session) error
	StoreRefreshToken(token *RefreshToken) error
	RevokeRefreshToken(tokenHash string) error
	GetRefreshToken(tokenHash string) (*RefreshToken, error)
}
