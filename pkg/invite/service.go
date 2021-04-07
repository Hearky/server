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

package invite

import (
	"context"
	"github.com/google/uuid"
	"github.com/hearky/server/pkg/domain"
	"time"
)

type service struct {
	inviteRepo  domain.InviteRepository
	meetingRepo domain.MeetingRepository
	userRepo    domain.UserRepository
}

func NewService(inviteRepository domain.InviteRepository, meetingRepository domain.MeetingRepository, userRepository domain.UserRepository) domain.InviteService {
	return &service{
		inviteRepo:  inviteRepository,
		meetingRepo: meetingRepository,
		userRepo:    userRepository,
	}
}

func (s *service) SendInvite(dto *domain.CreateInviteDto, uid string) error {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Check if the exact invite already exists
	_, err := s.inviteRepo.GetInviteByReceiverAndMeeting(ctx, dto.ReceiverID, dto.MeetingID)
	if err != nil && err != domain.ErrNotFound {
		return err
	} else if err != domain.ErrNotFound {
		return domain.ErrInviteExists
	}

	// Check if meeting exists
	m, err := s.meetingRepo.GetMeetingByID(ctx, dto.MeetingID)
	if err != nil {
		return err
	}

	// Check if user is organizer
	if !m.IsOrganizer(uid) {
		return domain.ErrForbidden
	}

	// Create new invite
	id := uuid.New().String()
	i := &domain.Invite{
		ID:         id,
		SenderID:   uid,
		ReceiverID: dto.ReceiverID,
		MeetingID:  dto.MeetingID,
		Timestamp:  time.Now(),
	}
	return s.inviteRepo.CreateInvite(ctx, i)
}

func (s *service) GetInvitesByReceiver(uid string) ([]*domain.Invite, error) {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Fetch invites
	i, err := s.inviteRepo.GetInvitesByReceiver(ctx, uid)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (s *service) GetInvitesByMeeting(mid string, uid string) ([]*domain.Invite, error) {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Fetch meeting
	m, err := s.meetingRepo.GetMeetingByID(ctx, mid)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if !m.IsParticipant(uid) {
		return nil, domain.ErrForbidden
	}

	// Fetch invites
	i, err := s.inviteRepo.GetInvitesByMeeting(ctx, mid)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (s *service) GetInvitesByReceiverCount(uid string) (int64, error) {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Fetch invites
	c, err := s.inviteRepo.GetInvitesByReceiverCount(ctx, uid)
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (s *service) GetInvitesByMeetingCount(mid string, uid string) (int64, error) {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Fetch meeting
	m, err := s.meetingRepo.GetMeetingByID(ctx, mid)
	if err != nil {
		return 0, err
	}

	// Check permissions
	if !m.IsParticipant(uid) {
		return 0, domain.ErrForbidden
	}

	// Fetch invites
	c, err := s.inviteRepo.GetInvitesByMeetingCount(ctx, mid)
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (s *service) AcceptInvite(id string, uid string) error {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Check if invite exists
	i, err := s.inviteRepo.GetInviteByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if the current user is the receiver
	if i.ReceiverID != uid {
		return domain.ErrForbidden
	}

	// Fetch meeting, if it does not exist -> delete invite
	m, err := s.meetingRepo.GetMeetingByID(ctx, i.MeetingID)
	if err != nil {
		return s.inviteRepo.DeleteInvite(ctx, id)
	}

	// Check if user still exists, if not -> delete invite
	_, err = s.userRepo.GetUserByID(ctx, i.ReceiverID)
	if err != nil {
		return s.inviteRepo.DeleteInvite(ctx, id)
	}

	// Add invited user as a participant to the meeting
	m.AddParticipant(i.ReceiverID)
	err = s.meetingRepo.SaveMeeting(ctx, m)
	if err != nil {
		return err
	}

	// Delete invite
	return s.inviteRepo.DeleteInvite(ctx, id)
}

func (s *service) DeleteInvite(id string, uid string) error {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Check if invite exists
	i, err := s.inviteRepo.GetInviteByID(ctx, id)
	if err != nil {
		return err
	}

	// Fetch meeting, if it does not exist -> delete invite
	m, err := s.meetingRepo.GetMeetingByID(ctx, i.MeetingID)
	if err != nil {
		return s.inviteRepo.DeleteInvite(ctx, id)
	}

	// If sender is not an organizer anymore -> delete invite
	if !m.IsOrganizer(i.SenderID) {
		return s.inviteRepo.DeleteInvite(ctx, id)
	}

	// Check if current user is an organizer
	if !m.IsOrganizer(uid) && i.ReceiverID != uid {
		return domain.ErrForbidden
	}

	// Delete invite
	return s.inviteRepo.DeleteInvite(ctx, id)
}
