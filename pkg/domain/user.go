package domain

import "context"

type CreateUserDto struct {
	Username string `json:"username"`
}

type User struct {
	ID       string      `json:"id" bson:"_id"`
	Username string      `json:"username"`
	Upgrade  UserUpgrade `json:"upgrade"`
}

type UserUpgrade struct {
	ConcurrentMeetings int `json:"concurrent_meetings" bson:"concurrent_meetings"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, u *User) error
	GetUserByID(ctx context.Context, uid string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	SaveUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, uid string) error
}

type UserService interface {
	CreateUser(dto *CreateUserDto, uid string) error
	GetUser(id string, uid string) (*User, error)
	DeleteUser(id string, uid string) error
}
