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

package user

import (
	"context"
	"github.com/hearky/server/pkg/domain"
	"time"
)

type service struct {
	userRepo    domain.UserRepository
	meetingRepo domain.MeetingRepository
	inviteRepo  domain.InviteRepository
}

func NewService(userRepository domain.UserRepository, meetingRepository domain.MeetingRepository, inviteRepository domain.InviteRepository) domain.UserService {
	return &service{
		userRepo:    userRepository,
		meetingRepo: meetingRepository,
		inviteRepo:  inviteRepository,
	}
}

func (s *service) CreateUser(dto *domain.CreateUserDto, uid string) error {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Check if account already exists
	_, err := s.userRepo.GetUserByID(ctx, uid)
	if err != nil && err != domain.ErrNotFound {
		return err
	} else if err != domain.ErrNotFound {
		return domain.ErrUsernameExists
	}

	// Check if username already exists
	_, err = s.userRepo.GetUserByUsername(ctx, dto.Username)
	if err != nil && err != domain.ErrNotFound {
		return err
	} else if err != domain.ErrNotFound {
		return domain.ErrUsernameExists
	}

	// Create new user
	u := &domain.User{
		ID:       uid,
		Username: dto.Username,
		Upgrade: domain.UserUpgrade{
			ConcurrentMeetings: domain.MaxConcurrentMeetings,
		},
	}
	return s.userRepo.CreateUser(ctx, u)
}

func (s *service) GetUser(id string, uid string) (*domain.User, error) {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Fetch account
	u, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	//TODO: Implement User Flag Check -> Moderator, Admin etc...
	if u.ID != uid {
		return nil, domain.ErrForbidden
	}
	return u, nil
}

func (s *service) DeleteUser(id string, uid string) error {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 20*time.Second)
	defer ccl()

	// Fetch account
	u, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	//TODO: Implement User Flag Check -> Moderator, Admin etc...
	if u.ID != uid {
		return domain.ErrForbidden
	}

	// Fetch meetings
	ms, err := s.meetingRepo.GetMeetingsByUser(ctx, id)
	if err != nil {
		return err
	}

	// Check if user is owner of meetings
	for _, m := range ms {
		if m.OwnerID == id {
			return domain.ErrOwnerOfMeeting
		}
	}

	// Fetch invites and delete them
	is, err := s.inviteRepo.GetInvitesByReceiver(ctx, id)
	if err != nil {
		return err
	}
	for _, i := range is {
		_ = s.inviteRepo.DeleteInvite(ctx, i.ID)
	}

	return nil
}
