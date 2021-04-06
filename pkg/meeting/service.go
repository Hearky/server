package meeting

import (
	"context"
	"github.com/google/uuid"
	"github.com/hearky/server/pkg/domain"
)

type Service struct {
	meetingRepo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{
		meetingRepo: r,
	}
}

func (s *Service) CreateMeeting(ctx context.Context, dto *domain.CreateMeetingDto, ownerId string) (string, error) {
	id := uuid.New().String()
	m := &domain.Meeting{
		ID:    id,
		Name:  dto.Name,
		Owner: ownerId,
	}
	//TODO: Iterate through organizers and participants and invite them
	return id, s.meetingRepo.CreateMeeting(ctx, m)
}

func (s *Service) GetMeetingByID(ctx context.Context, id string) (*domain.Meeting, error) {
	return s.meetingRepo.GetMeetingByID(ctx, id)
}

func (s *Service) GetMeetingsByUserID(ctx context.Context, id string) ([]*domain.Meeting, error) {
	return s.meetingRepo.GetMeetingsByUserID(ctx, id)
}

func (s *Service) DeleteMeeting(ctx context.Context, id string) error {
	return s.meetingRepo.DeleteMeetingByID(ctx, id)
}