package user

import "context"

type User struct {
	ID           int64
	Name         string
	PasswordHash string
}

type Bar struct {
	NewName string
	NewPassword string
}

type UserRepository interface {
	CreateUser(ctx context.Context, u User) error
	GetUserById(ctx context.Context, uID int64) (User, error)
	GetUserByName(ctx context.Context, uName string) (User, error)
	UpdateUserById(ctx context.Context, uID int64, dataToChange Bar) error
	DeleteUserById(ctx context.Context, uID int64) error
}
