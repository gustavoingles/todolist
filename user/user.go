package user

import "context"

type User struct {
	ID           int64
	Name         string
	PasswordHash string
}

type CreateUserCommand struct {
	Name     string
	Password string
}

type UserRepository interface {
	CreateUser(ctx context.Context, cmd CreateUserCommand) error
	GetUserById(ctx context.Context, uID int64) (User, error)
	GetUserByName(ctx context.Context, uName string) (User, error)
	UpdateUserName(ctx context.Context, updateFn func(*User) error) error
	DeleteUserById(ctx context.Context, uID int64) error
}
