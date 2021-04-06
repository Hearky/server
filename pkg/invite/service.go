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
	// Define Timeout
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

func (s *service) GetInvitesByReceiver(uid string) error {
	panic("implement me")
}

func (s *service) AcceptInvite(id string) error {
	// Define Timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Check if invite exists
	i, err := s.inviteRepo.GetInviteByID(ctx, id)
	if err != nil {
		return err
	}

	// Fetch meeting
	m, err := s.meetingRepo.GetMeetingByID(ctx, i.MeetingID)
	if err != nil {
		return err
	}

	// Check if user still exists
	_, err = s.userRepo.GetUserByID(ctx, i.ReceiverID)
	if err != nil {
		return err
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

func (s *service) RejectInvite(id string) error {
	// Define Timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Check if invite exists
	_, err := s.inviteRepo.GetInviteByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete invite
	return s.inviteRepo.DeleteInvite(ctx, id)
}
