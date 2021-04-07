/*
 * Hearky Server
 * Copyright (C) 2021 Hearky
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
