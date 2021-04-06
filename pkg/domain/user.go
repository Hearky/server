package domain

import "context"

type CreateUserDto struct {
	Username string `json:"username"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UserRepository interface {
	GetUserByID(ctx context.Context, uid string) (*User, error)
}
