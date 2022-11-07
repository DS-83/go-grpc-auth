package auth

import (
	context "context"
	"example-grpc-auth/models"
)

// Users storage interface
type UserRepo interface {
	CreateUser(c context.Context, u string, p string) error
	GetUser(c context.Context, u string, p string) (*models.User, error)
	UpdateUser(c context.Context, f *models.User, u *models.User) (*models.User, error)
	DeleteUser(context.Context, *models.User) error
}

// Tokens storage interface
type TokenRepo interface {
	RevokeToken(c context.Context, t string) error
	IsRevoked(c context.Context, t string) (bool, error)
}
