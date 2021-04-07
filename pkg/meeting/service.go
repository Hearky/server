package meeting

import (
	"context"
	"github.com/google/uuid"
	"github.com/hearky/server/pkg/domain"
	"time"
)

type service struct {
	meetingRepo domain.MeetingRepository
	inviteRepo  domain.InviteRepository
	userRepo    domain.UserRepository
}

func NewService(meetingRepository domain.MeetingRepository, inviteRepository domain.InviteRepository, userRepository domain.UserRepository) domain.MeetingService {
	return &service{
		meetingRepo: meetingRepository,
		inviteRepo:  inviteRepository,
		userRepo:    userRepository,
	}
}

func (s *service) CreateMeeting(dto *domain.CreateMeetingDto, uid string) (string, error) {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 20*time.Second)
	defer ccl()

	u, err := s.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		return "", err
	}

	ms, err := s.meetingRepo.GetMeetingsByUserCount(ctx, uid)
	if err != nil {
		return "", err
	}

	if ms+1 > int64(u.Upgrade.ConcurrentMeetings) {
		return "", domain.ErrTooManyMeetings
	}

	// Create meeting
	id := uuid.New().String()
	m := &domain.Meeting{
		ID:      id,
		Name:    dto.Name,
		OwnerID: uid,
		Upgrade: domain.MeetingUpgrade{
			Participants:      domain.MaxParticipantsFree,
			ConcurrentInvites: domain.MaxConcurrentInvitesFree,
		},
	}
	err = s.meetingRepo.CreateMeeting(ctx, m)
	if err != nil {
		return "", err
	}

	// Invite all participants
	if len(dto.Participants) > domain.MaxConcurrentInvitesFree {
		dto.Participants = dto.Participants[0 : domain.MaxConcurrentInvitesFree-1]
	}
	for _, p := range dto.Participants {
		// TODO: Check if they are friends
		_, err := s.userRepo.GetUserByID(ctx, p)
		if err != nil {
			continue
		}
		_ = s.inviteRepo.CreateInvite(ctx, &domain.Invite{
			ID:         uuid.New().String(),
			SenderID:   uid,
			ReceiverID: p,
			MeetingID:  m.ID,
			Timestamp:  time.Now(),
		})
	}

	return id, nil
}

func (s *service) GetMeetingByID(mid string, uid string) (*domain.Meeting, error) {
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

	return m, nil
}

func (s *service) GetMeetingsByUser(uid string) ([]*domain.Meeting, error) {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Fetch meetings
	m, err := s.meetingRepo.GetMeetingsByUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *service) GetMeetingsByUserCount(uid string) (int64, error) {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Fetch meetings count
	c, err := s.meetingRepo.GetMeetingsByUserCount(ctx, uid)
	if err != nil {
		return -1, err
	}

	return c, nil
}

func (s *service) DeleteMeeting(mid string, uid string) error {
	// Define timeout
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()

	// Fetch meeting
	m, err := s.meetingRepo.GetMeetingByID(ctx, mid)
	if err != nil {
		return err
	}

	// Check permissions
	if !m.IsOwner(uid) {
		return domain.ErrForbidden
	}

	// Delete all invites
	is, err := s.inviteRepo.GetInvitesByMeeting(ctx, mid)
	for _, i := range is {
		_ = s.inviteRepo.DeleteInvite(ctx, i.ID)
	}

	// Delete meeting
	return s.meetingRepo.DeleteMeeting(ctx, mid)
}
